package digest

import (
	"context"
	"fmt"
	didstate "github.com/ProtoconNet/mitum-credential/state"
	crcystate "github.com/ProtoconNet/mitum-currency/v3/state"
	nftstate "github.com/ProtoconNet/mitum-nft/state"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"strconv"
	"sync"
	"time"

	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-currency/v3/digest/isaac"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stateextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	mitumutil "github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/fixedtree"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var bulkWriteLimit = 500

type BlockSession struct {
	sync.RWMutex
	block                      mitumbase.BlockMap
	ops                        []mitumbase.Operation
	opsTree                    fixedtree.Tree
	sts                        []mitumbase.State
	st                         *currencydigest.Database
	proposal                   mitumbase.ProposalSignFact
	opsTreeNodes               map[string]mitumbase.OperationFixedtreeNode
	blockModels                []mongo.WriteModel
	operationModels            []mongo.WriteModel
	accountModels              []mongo.WriteModel
	balanceModels              []mongo.WriteModel
	currencyModels             []mongo.WriteModel
	contractAccountModels      []mongo.WriteModel
	nftCollectionModels        []mongo.WriteModel
	nftModels                  []mongo.WriteModel
	nftBoxModels               []mongo.WriteModel
	nftOperatorModels          []mongo.WriteModel
	didIssuerModels            []mongo.WriteModel
	didCredentialModels        []mongo.WriteModel
	didHolderDIDModels         []mongo.WriteModel
	didTemplateModels          []mongo.WriteModel
	timestampModels            []mongo.WriteModel
	tokenModels                []mongo.WriteModel
	tokenBalanceModels         []mongo.WriteModel
	pointModels                []mongo.WriteModel
	pointBalanceModels         []mongo.WriteModel
	daoDesignModels            []mongo.WriteModel
	daoProposalModels          []mongo.WriteModel
	daoDelegatorsModels        []mongo.WriteModel
	daoVotersModels            []mongo.WriteModel
	daoVotingPowerBoxModels    []mongo.WriteModel
	storageModels              []mongo.WriteModel
	storageDataModels          []mongo.WriteModel
	prescriptionModels         []mongo.WriteModel
	prescriptionInfoDataModels []mongo.WriteModel
	didRegistryModels          []mongo.WriteModel
	didDataModels              []mongo.WriteModel
	didDocumentModels          []mongo.WriteModel
	dMileModels                []mongo.WriteModel
	dMileDataModels            []mongo.WriteModel
	statesValue                *sync.Map
	balanceAddressList         []string
	nftMap                     map[string]struct{}
	credentialMap              map[string]struct{}
	buildinfo                  string
}

func NewBlockSession(
	st *currencydigest.Database,
	blk mitumbase.BlockMap,
	ops []mitumbase.Operation,
	opstree fixedtree.Tree,
	sts []mitumbase.State,
	proposal mitumbase.ProposalSignFact,
	vs string,
) (*BlockSession, error) {
	if st.Readonly() {
		return nil, errors.Errorf("readonly mode")
	}

	nst, err := st.New()
	if err != nil {
		return nil, err
	}

	return &BlockSession{
		st:            nst,
		block:         blk,
		ops:           ops,
		opsTree:       opstree,
		sts:           sts,
		proposal:      proposal,
		statesValue:   &sync.Map{},
		nftMap:        map[string]struct{}{},
		credentialMap: map[string]struct{}{},
		buildinfo:     vs,
	}, nil
}

func (bs *BlockSession) Prepare() error {
	bs.Lock()
	defer bs.Unlock()
	if err := bs.prepareOperationsTree(); err != nil {
		return err
	}
	if err := bs.prepareBlock(); err != nil {
		return err
	}
	if err := bs.prepareOperations(); err != nil {
		return err
	}
	if err := bs.prepareCurrencies(); err != nil {
		return err
	}
	if err := bs.prepareNFTs(); err != nil {
		return err
	}
	if err := bs.prepareDIDCredential(); err != nil {
		return err
	}
	if err := bs.prepareTimeStamps(); err != nil {
		return err
	}
	if err := bs.prepareToken(); err != nil {
		return err
	}
	if err := bs.preparePoint(); err != nil {
		return err
	}
	if err := bs.prepareDAO(); err != nil {
		return err
	}
	if err := bs.prepareStorage(); err != nil {
		return err
	}
	if err := bs.preparePrescription(); err != nil {
		return err
	}
	if err := bs.prepareDIDRegistry(); err != nil {
		return err
	}
	if err := bs.prepareDmile(); err != nil {
		return err
	}

	return bs.prepareAccounts()
}

func (bs *BlockSession) Commit(ctx context.Context) error {
	bs.Lock()
	defer bs.Unlock()

	started := time.Now()
	defer func() {
		bs.statesValue.Store("commit", time.Since(started))

		_ = bs.close()
	}()
	wc := writeconcern.Majority()
	sessionOpts := options.Session().SetCausalConsistency(true)

	session, err := bs.st.MongoClient().MongoClient().StartSession(sessionOpts)
	if err != nil {
		return err
	}
	defer session.EndSession(context.TODO())

	txnOpts := options.Transaction().
		SetReadConcern(readconcern.Snapshot()).
		SetWriteConcern(wc).
		SetReadPreference(readpref.Primary())

	err = mongo.WithSession(ctx, session, func(txnCtx mongo.SessionContext) error {
		// 트랜잭션 시작
		if err := session.StartTransaction(txnOpts); err != nil {
			return err
		}

		if err := bs.writeModels(txnCtx, defaultColNameBlock, bs.blockModels); err != nil {
			session.AbortTransaction(txnCtx)
			return err
		}

		if len(bs.operationModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameOperation, bs.operationModels); err != nil {
				session.AbortTransaction(txnCtx)
				return err
			}
		}

		if len(bs.currencyModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameCurrency, bs.currencyModels); err != nil {
				session.AbortTransaction(txnCtx)
				return err
			}
		}

		if len(bs.accountModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameAccount, bs.accountModels); err != nil {
				session.AbortTransaction(txnCtx)
				return err
			}
		}

		if len(bs.contractAccountModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameContractAccount, bs.contractAccountModels); err != nil {
				session.AbortTransaction(txnCtx)
				return err
			}
		}

		if len(bs.balanceModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameBalance, bs.balanceModels); err != nil {
				return err
			}
		}

		if len(bs.nftCollectionModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameNFTCollection, bs.nftCollectionModels); err != nil {
				return err
			}
		}

		if len(bs.nftModels) > 0 {
			for key := range bs.nftMap {
				parsedKey, err := crcystate.ParseStateKey(key, nftstate.NFTPrefix, 4)
				if err != nil {
					return err
				}
				i, _ := strconv.ParseInt(parsedKey[2], 10, 64)
				err = bs.st.CleanByHeightColName(
					txnCtx,
					bs.block.Manifest().Height(),
					defaultColNameNFT,
					bson.D{{"contract", parsedKey[1]}},
					bson.D{{"nft_idx", i}},
				)
				if err != nil {
					return err
				}
			}

			if err := bs.writeModels(txnCtx, defaultColNameNFT, bs.nftModels); err != nil {
				return err
			}
		}

		if len(bs.nftOperatorModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameNFTOperator, bs.nftOperatorModels); err != nil {
				return err
			}
		}

		if len(bs.nftBoxModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameNFT, bs.nftBoxModels); err != nil {
				return err
			}
		}

		if len(bs.balanceModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameBalance, bs.balanceModels); err != nil {
				return err
			}
		}

		if len(bs.didIssuerModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameDIDCredentialService, bs.didIssuerModels); err != nil {
				return err
			}
		}

		if len(bs.didCredentialModels) > 0 {
			for key := range bs.credentialMap {
				parsedKey, err := crcystate.ParseStateKey(key, didstate.CredentialPrefix, 5)
				if err != nil {
					return err
				}
				err = bs.st.CleanByHeightColName(
					txnCtx,
					bs.block.Manifest().Height(),
					defaultColNameDIDCredential,
					bson.D{{"contract", parsedKey[1]}},
					bson.D{{"template", parsedKey[2]}},
					bson.D{{"credential_id", parsedKey[3]}},
				)
				if err != nil {
					return err
				}
			}

			if err := bs.writeModels(txnCtx, defaultColNameDIDCredential, bs.didCredentialModels); err != nil {
				return err
			}
		}

		if len(bs.didHolderDIDModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameHolder, bs.didHolderDIDModels); err != nil {
				return err
			}
		}

		if len(bs.didTemplateModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameTemplate, bs.didTemplateModels); err != nil {
				return err
			}
		}

		if len(bs.timestampModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameTimeStamp, bs.timestampModels); err != nil {
				return err
			}
		}

		if len(bs.tokenModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameToken, bs.tokenModels); err != nil {
				return err
			}
		}

		if len(bs.tokenBalanceModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameTokenBalance, bs.tokenBalanceModels); err != nil {
				return err
			}
		}

		if len(bs.pointModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNamePoint, bs.pointModels); err != nil {
				return err
			}
		}

		if len(bs.pointBalanceModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNamePointBalance, bs.pointBalanceModels); err != nil {
				return err
			}
		}

		if len(bs.daoDesignModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameDAO, bs.daoDesignModels); err != nil {
				return err
			}
		}

		if len(bs.daoProposalModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameDAOProposal, bs.daoProposalModels); err != nil {
				return err
			}
		}

		if len(bs.daoDelegatorsModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameDAODelegators, bs.daoDelegatorsModels); err != nil {
				return err
			}
		}

		if len(bs.daoVotersModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameDAOVoters, bs.daoVotersModels); err != nil {
				return err
			}
		}

		if len(bs.daoVotingPowerBoxModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameDAOVotingPowerBox, bs.daoVotingPowerBoxModels); err != nil {
				return err
			}
		}

		if len(bs.storageModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameStorage, bs.storageModels); err != nil {
				return err
			}
		}

		if len(bs.storageDataModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameStorageData, bs.storageDataModels); err != nil {
				return err
			}
		}

		if len(bs.prescriptionModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNamePrescription, bs.prescriptionModels); err != nil {
				return err
			}
		}

		if len(bs.prescriptionInfoDataModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNamePrescriptionInfo, bs.prescriptionInfoDataModels); err != nil {
				return err
			}
		}

		if len(bs.didRegistryModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameDIDRegistry, bs.didRegistryModels); err != nil {
				return err
			}
		}

		if len(bs.didDataModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameDIDData, bs.didDataModels); err != nil {
				return err
			}
		}

		if len(bs.didDocumentModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameDIDDocument, bs.didDocumentModels); err != nil {
				return err
			}
		}

		if len(bs.dMileModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameDmile, bs.dMileModels); err != nil {
				return err
			}
		}

		if len(bs.dMileDataModels) > 0 {
			if err := bs.writeModels(txnCtx, defaultColNameDmileData, bs.dMileDataModels); err != nil {
				return err
			}
		}

		if err := session.CommitTransaction(txnCtx); err != nil {
			return err
		}

		//time.Sleep(1000 * time.Millisecond)

		return nil
	})

	//_, err := bs.st.MongoClient().WithSession(func(txnCtx mongo.SessionContext, collection func(string) *mongo.Collection) (interface{}, error) {
	//	if err := bs.writeModels(txnCtx, defaultColNameBlock, bs.blockModels); err != nil {
	//		return nil, err
	//	}
	//
	//	if len(bs.operationModels) > 0 {
	//		if err := bs.writeModels(txnCtx, defaultColNameOperation, bs.operationModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	if len(bs.currencyModels) > 0 {
	//		if err := bs.writeModels(txnCtx, defaultColNameCurrency, bs.currencyModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	if len(bs.accountModels) > 0 {
	//		if err := bs.writeModels(txnCtx, defaultColNameAccount, bs.accountModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	if len(bs.contractAccountModels) > 0 {
	//		if err := bs.writeModels(txnCtx, defaultColNameContractAccount, bs.contractAccountModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	if len(bs.balanceModels) > 0 {
	//		if err := bs.writeModels(txnCtx, defaultColNameBalance, bs.balanceModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	if len(bs.nftCollectionModels) > 0 {
	//		if err := bs.writeModels(txnCtx, defaultColNameNFTCollection, bs.nftCollectionModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	if len(bs.nftModels) > 0 {
	//		for key := range bs.nftMap {
	//			parsedKey, err := crcystate.ParseStateKey(key, nftstate.NFTPrefix, 4)
	//			if err != nil {
	//				return nil, err
	//			}
	//			i, _ := strconv.ParseInt(parsedKey[2], 10, 64)
	//			err = bs.st.CleanByHeightColName(
	//				txnCtx,
	//				bs.block.Manifest().Height(),
	//				defaultColNameNFT,
	//				bson.D{{"contract", parsedKey[1]}},
	//				bson.D{{"nft_idx", i}},
	//			)
	//			if err != nil {
	//				return nil, err
	//			}
	//		}
	//
	//		if err := bs.writeModels(txnCtx, defaultColNameNFT, bs.nftModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	if len(bs.nftOperatorModels) > 0 {
	//		if err := bs.writeModels(txnCtx, defaultColNameNFTOperator, bs.nftOperatorModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	if len(bs.nftBoxModels) > 0 {
	//		if err := bs.writeModels(txnCtx, defaultColNameNFT, bs.nftBoxModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	if len(bs.balanceModels) > 0 {
	//		if err := bs.writeModels(txnCtx, defaultColNameBalance, bs.balanceModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	if len(bs.didIssuerModels) > 0 {
	//		if err := bs.writeModels(txnCtx, defaultColNameDIDCredentialService, bs.didIssuerModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	if len(bs.didCredentialModels) > 0 {
	//		for key := range bs.credentialMap {
	//			parsedKey, err := crcystate.ParseStateKey(key, didstate.CredentialPrefix, 5)
	//			if err != nil {
	//				return nil, err
	//			}
	//			err = bs.st.CleanByHeightColName(
	//				txnCtx,
	//				bs.block.Manifest().Height(),
	//				defaultColNameDIDCredential,
	//				bson.D{{"contract", parsedKey[1]}},
	//				bson.D{{"template", parsedKey[2]}},
	//				bson.D{{"credential_id", parsedKey[3]}},
	//			)
	//			if err != nil {
	//				return nil, err
	//			}
	//		}
	//
	//		if err := bs.writeModels(txnCtx, defaultColNameDIDCredential, bs.didCredentialModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	if len(bs.didHolderDIDModels) > 0 {
	//		if err := bs.writeModels(txnCtx, defaultColNameHolder, bs.didHolderDIDModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	if len(bs.didTemplateModels) > 0 {
	//		if err := bs.writeModels(txnCtx, defaultColNameTemplate, bs.didTemplateModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	if len(bs.timestampModels) > 0 {
	//		if err := bs.writeModels(txnCtx, defaultColNameTimeStamp, bs.timestampModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	if len(bs.tokenModels) > 0 {
	//		if err := bs.writeModels(txnCtx, defaultColNameToken, bs.tokenModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	if len(bs.tokenBalanceModels) > 0 {
	//		if err := bs.writeModels(txnCtx, defaultColNameTokenBalance, bs.tokenBalanceModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	if len(bs.pointModels) > 0 {
	//		if err := bs.writeModels(txnCtx, defaultColNamePoint, bs.pointModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	if len(bs.pointBalanceModels) > 0 {
	//		if err := bs.writeModels(txnCtx, defaultColNamePointBalance, bs.pointBalanceModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	if len(bs.daoDesignModels) > 0 {
	//		if err := bs.writeModels(txnCtx, defaultColNameDAO, bs.daoDesignModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	if len(bs.daoProposalModels) > 0 {
	//		if err := bs.writeModels(txnCtx, defaultColNameDAOProposal, bs.daoProposalModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	if len(bs.daoDelegatorsModels) > 0 {
	//		if err := bs.writeModels(txnCtx, defaultColNameDAODelegators, bs.daoDelegatorsModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	if len(bs.daoVotersModels) > 0 {
	//		if err := bs.writeModels(txnCtx, defaultColNameDAOVoters, bs.daoVotersModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	if len(bs.daoVotingPowerBoxModels) > 0 {
	//		if err := bs.writeModels(txnCtx, defaultColNameDAOVotingPowerBox, bs.daoVotingPowerBoxModels); err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	return nil, nil
	//})
	time.Sleep(1000 * time.Millisecond)

	return err
}

func (bs *BlockSession) Close() error {
	bs.Lock()
	defer bs.Unlock()

	return bs.close()
}

func (bs *BlockSession) prepareOperationsTree() error {
	nodes := map[string]mitumbase.OperationFixedtreeNode{}

	if err := bs.opsTree.Traverse(func(_ uint64, no fixedtree.Node) (bool, error) {
		nno := no.(mitumbase.OperationFixedtreeNode)

		if nno.Reason() == nil {
			nodes[nno.Key()] = nno
		} else {
			nodes[nno.Key()[:len(nno.Key())-1]] = nno
		}

		return true, nil
	}); err != nil {
		return err
	}

	bs.opsTreeNodes = nodes

	return nil
}

func (bs *BlockSession) prepareBlock() error {
	if bs.block == nil {
		return nil
	}

	bs.blockModels = make([]mongo.WriteModel, 1)

	manifest := isaac.NewManifest(
		bs.block.Manifest().Height(),
		bs.block.Manifest().Previous(),
		bs.block.Manifest().Proposal(),
		bs.block.Manifest().OperationsTree(),
		bs.block.Manifest().StatesTree(),
		bs.block.Manifest().Suffrage(),
		bs.block.Manifest().ProposedAt(),
	)

	doc, err := currencydigest.NewManifestDoc(manifest, bs.st.Encoder(), bs.block.Manifest().Height(), bs.ops, bs.block.SignedAt(), bs.proposal.ProposalFact().Proposer(), bs.proposal.ProposalFact().Point().Round(), bs.buildinfo)
	if err != nil {
		return err
	}
	bs.blockModels[0] = mongo.NewInsertOneModel().SetDocument(doc)

	return nil
}

func (bs *BlockSession) prepareOperations() error {
	if len(bs.ops) < 1 {
		return nil
	}

	node := func(h mitumutil.Hash) (bool, bool, mitumbase.OperationProcessReasonError) {
		no, found := bs.opsTreeNodes[h.String()]
		if !found {
			return false, false, nil
		}

		return true, no.InState(), no.Reason()
	}

	bs.operationModels = make([]mongo.WriteModel, len(bs.ops))

	for i := range bs.ops {
		op := bs.ops[i]

		var doc currencydigest.OperationDoc
		switch found, inState, reason := node(op.Fact().Hash()); {
		case !found:
			return mitumutil.ErrNotFound.Errorf("operation, %v in operations tree", op.Fact().Hash().String())
		default:
			var reasonMsg string
			switch {
			case reason == nil:
				reasonMsg = ""
			default:
				reasonMsg = reason.Msg()
			}
			d, err := currencydigest.NewOperationDoc(
				op,
				bs.st.Encoder(),
				bs.block.Manifest().Height(),
				bs.block.SignedAt(),
				inState,
				reasonMsg,
				uint64(i),
			)
			if err != nil {
				return err
			}
			doc = d
		}

		bs.operationModels[i] = mongo.NewInsertOneModel().SetDocument(doc)
	}

	return nil
}

func (bs *BlockSession) prepareAccounts() error {
	if len(bs.sts) < 1 {
		return nil
	}

	var accountModels []mongo.WriteModel
	var balanceModels []mongo.WriteModel
	var contractAccountModels []mongo.WriteModel
	for i := range bs.sts {
		st := bs.sts[i]

		switch {
		case statecurrency.IsAccountStateKey(st.Key()):
			j, err := bs.handleAccountState(st)
			if err != nil {
				return err
			}
			accountModels = append(accountModels, j...)
		case statecurrency.IsBalanceStateKey(st.Key()):
			j, address, err := bs.handleBalanceState(st)
			if err != nil {
				return err
			}
			balanceModels = append(balanceModels, j...)
			bs.balanceAddressList = append(bs.balanceAddressList, address)
		case stateextension.IsStateContractAccountKey(st.Key()):
			j, err := bs.handleContractAccountState(st)
			if err != nil {
				return err
			}
			contractAccountModels = append(contractAccountModels, j...)
		default:
			continue
		}
	}

	bs.accountModels = accountModels
	bs.contractAccountModels = contractAccountModels
	bs.balanceModels = balanceModels

	return nil
}

func (bs *BlockSession) prepareCurrencies() error {
	if len(bs.sts) < 1 {
		return nil
	}

	var currencyModels []mongo.WriteModel
	for i := range bs.sts {
		st := bs.sts[i]
		switch {
		case statecurrency.IsDesignStateKey(st.Key()):
			j, err := bs.handleCurrencyState(st)
			if err != nil {
				return err
			}
			currencyModels = append(currencyModels, j...)
		default:
			continue
		}
	}

	bs.currencyModels = currencyModels

	return nil
}

func (bs *BlockSession) writeModels(ctx context.Context, col string, models []mongo.WriteModel) error {
	started := time.Now()
	defer func() {
		bs.statesValue.Store(fmt.Sprintf("write-models-%s", col), time.Since(started))
	}()

	n := len(models)
	if n < 1 {
		return nil
	} else if n <= bulkWriteLimit {
		return bs.writeModelsChunk(ctx, col, models)
	}

	z := n / bulkWriteLimit
	if n%bulkWriteLimit != 0 {
		z++
	}

	for i := 0; i < z; i++ {
		s := i * bulkWriteLimit
		e := s + bulkWriteLimit
		if e > n {
			e = n
		}

		if err := bs.writeModelsChunk(ctx, col, models[s:e]); err != nil {
			return err
		}
	}

	return nil
}

func (bs *BlockSession) writeModelsChunk(ctx context.Context, col string, models []mongo.WriteModel) error {
	opts := options.BulkWrite().SetOrdered(true)
	if res, err := bs.st.MongoClient().Collection(col).BulkWrite(ctx, models, opts); err != nil {
		return err
	} else if res != nil && res.InsertedCount < 1 {
		return errors.Errorf("not inserted to %s", col)
	}

	return nil
}

func (bs *BlockSession) close() error {
	bs.block = nil
	bs.ops = nil
	bs.opsTree = fixedtree.EmptyTree()
	bs.sts = nil
	bs.proposal = nil
	bs.opsTreeNodes = nil
	bs.blockModels = nil
	bs.operationModels = nil
	bs.currencyModels = nil
	bs.accountModels = nil
	bs.balanceModels = nil
	bs.contractAccountModels = nil
	bs.nftCollectionModels = nil
	bs.nftModels = nil
	bs.nftOperatorModels = nil
	bs.didIssuerModels = nil
	bs.didCredentialModels = nil
	bs.didHolderDIDModels = nil
	bs.didTemplateModels = nil
	bs.timestampModels = nil
	bs.tokenModels = nil
	bs.tokenBalanceModels = nil
	bs.pointModels = nil
	bs.pointBalanceModels = nil
	bs.storageModels = nil
	bs.storageDataModels = nil
	bs.prescriptionModels = nil
	bs.prescriptionInfoDataModels = nil
	bs.didRegistryModels = nil
	bs.didDataModels = nil
	bs.didDocumentModels = nil
	bs.dMileModels = nil
	bs.dMileDataModels = nil
	bs.contractAccountModels = nil
	bs.nftMap = nil
	bs.credentialMap = nil

	return bs.st.Close()
}
