package digest

import (
	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	utilc "github.com/ProtoconNet/mitum-currency/v3/digest/util"
	"github.com/ProtoconNet/mitum-did-registry/state"
	"github.com/ProtoconNet/mitum-did-registry/types"
	"github.com/ProtoconNet/mitum2/base"
	utilm "github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DIDDesign(st *cdigest.Database, contract string) (types.Design, base.State, error) {
	filter := utilc.NewBSONFilter("contract", contract)
	q := filter.D()

	opt := options.FindOne().SetSort(
		utilc.NewBSONFilter("height", -1).D(),
	)
	var sta base.State
	if err := st.MongoClient().GetByFilter(
		defaultColNameDIDRegistry,
		q,
		func(res *mongo.SingleResult) error {
			i, err := cdigest.LoadState(res.Decode, st.Encoders())
			if err != nil {
				return err
			}
			sta = i
			return nil
		},
		opt,
	); err != nil {
		return types.Design{}, nil, utilm.ErrNotFound.WithMessage(err, "storage design by contract account %v", contract)
	}

	if sta != nil {
		de, err := state.GetDesignFromState(sta)
		if err != nil {
			return types.Design{}, nil, err
		}
		return de, sta, nil
	} else {
		return types.Design{}, nil, errors.Errorf("state is nil")
	}
}

func DIDData(db *cdigest.Database, contract, key string) (*types.Data, base.State, error) {
	filter := utilc.NewBSONFilter("contract", contract)
	filter = filter.Add("publicKey", key)
	q := filter.D()

	opt := options.FindOne().SetSort(
		utilc.NewBSONFilter("height", -1).D(),
	)
	var data *types.Data
	var sta base.State
	var err error
	if err := db.MongoClient().GetByFilter(
		defaultColNameDIDData,
		q,
		func(res *mongo.SingleResult) error {
			sta, err = cdigest.LoadState(res.Decode, db.Encoders())
			if err != nil {
				return err
			}
			d, err := state.GetDataFromState(sta)
			if err != nil {
				return err
			}
			data = &d
			return nil
		},
		opt,
	); err != nil {
		return nil, nil, utilm.ErrNotFound.WithMessage(
			err, "DID data for public key %s in contract account %s", key, contract)
	}

	if data != nil {
		return data, sta, nil
	} else {
		return nil, nil, errors.Errorf("data is nil")
	}
}

func DIDDocument(db *cdigest.Database, contract, key string) (*types.Document, base.State, error) {
	filter := utilc.NewBSONFilter("contract", contract)
	filter = filter.Add("did", key)
	q := filter.D()

	opt := options.FindOne().SetSort(
		utilc.NewBSONFilter("height", -1).D(),
	)
	var document *types.Document
	var sta base.State
	var err error
	if err := db.MongoClient().GetByFilter(
		defaultColNameDIDDocument,
		q,
		func(res *mongo.SingleResult) error {
			sta, err = cdigest.LoadState(res.Decode, db.Encoders())
			if err != nil {
				return err
			}
			d, err := state.GetDocumentFromState(sta)
			if err != nil {
				return err
			}
			document = &d
			return nil
		},
		opt,
	); err != nil {
		return nil, nil, utilm.ErrNotFound.WithMessage(
			err, "DID document for DID %s in contract account %s", key, contract)
	}

	if document != nil {
		return document, sta, nil
	} else {
		return nil, nil, errors.Errorf("document is nil")
	}
}
