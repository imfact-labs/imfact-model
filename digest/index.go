package digest

import (
	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var indexPrefix = "mitum_digest_"

var accountIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "address", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "account"),
	},
	{
		Keys: bson.D{bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "account_height"),
	},
	{
		Keys: bson.D{bson.E{Key: "pubs", Value: 1}, bson.E{Key: "height", Value: 1}, bson.E{Key: "address", Value: 1}},
		Options: options.Index().
			SetName(indexPrefix + "account_publiskeys"),
	},
}

var balanceIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{bson.E{Key: "address", Value: 1}, bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "balance"),
	},
	{
		Keys: bson.D{
			bson.E{Key: "address", Value: 1},
			bson.E{Key: "currency", Value: 1},
			bson.E{Key: "height", Value: -1},
		},
		Options: options.Index().
			SetName(indexPrefix + "balance_currency"),
	},
	//{
	//	Keys: bson.D{bson.E{Key: "height", Value: -1}},
	//	Options: options.Index().
	//		SetName(indexPrefix + "balance_height"),
	//},
}

var operationIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{bson.E{Key: "addresses", Value: 1}, bson.E{Key: "height", Value: 1}, bson.E{Key: "index", Value: 1}},
		Options: options.Index().
			SetName(indexPrefix + "account_operation"),
	},
	{
		Keys: bson.D{bson.E{Key: "height", Value: 1}, bson.E{Key: "index", Value: 1}},
		Options: options.Index().
			SetName(indexPrefix + "operation"),
	},
	{
		Keys: bson.D{bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "operation_height"),
	},
}

var timestampIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "height", Value: -1},
			bson.E{Key: "isItem", Value: 1},
			bson.E{Key: "project_id", Value: 1},
			bson.E{Key: "timestamp_idx", Value: 1}},
		Options: options.Index().
			SetName(indexPrefix + "timestamp_idx_contract_height_isItem_projectID"),
	},
}

var nftCollectionIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "nft_collection_contract_height"),
	},
}

var nftIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "nft_idx", Value: 1},
			bson.E{Key: "height", Value: -1},
			bson.E{Key: "istoken", Value: 1},
		},
		Options: options.Index().
			SetName(indexPrefix + "nft_idx_contract_height_istoken"),
	},
	{
		Keys: bson.D{bson.E{Key: "facthash", Value: 1}},
		Options: options.Index().
			SetName(indexPrefix + "nft_facthash"),
	},
}

var nftOperatorIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "address", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "nft_operator_contract_address_height"),
	},
}

var didCredentialServiceIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "credential_service_contract_height"),
	},
}

var didCredentialIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "template", Value: 1},
			bson.E{Key: "credential_id", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "credential_id_contract_template_height"),
	},
}

var didCredentialHolderIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "holder", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "credential_holder_contract_height"),
	},
}

var didCredentialTemplateIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "template", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "credential_template_contract_height"),
	},
}

var daoServiceIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "dao_service_contract_height"),
	},
}

var daoProposalIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "proposal_id", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "dao_proposal_contract_height"),
	},
}

var daoDelegatorsIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "proposal_id", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "dao_approved_contract_proposalID_height"),
	},
}

var daoVotersIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "proposal_id", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "dao_voter_contract_proposalID_height"),
	},
}

var daoVotingPowerBoxIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "proposal_id", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "dao_voting_power_contract_proposalID_height"),
	},
}

var tokenServiceIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "token_service_contract_height"),
	},
}

var tokenBalanceIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "address", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "token_balance_contract_address_height"),
	},
}

var pointServiceIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "point_service_contract_height"),
	},
}

var pointBalanceIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "address", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "point_balance_contract_address_height"),
	},
}

var storageDataServiceIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "storage_data_service_contract_height"),
	},
}

var storageDataIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "data_key", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "storage_data_contract_datakey_height"),
	},
}

var didRegistryIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "did_registry_contract_height"),
	},
}

var didRegistryDataIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "publicKey", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "did_registry_data_contract_publicKey_height"),
	},
}

var didRegistryDocumentIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "did", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(indexPrefix + "did_registry_document_contract_did_height"),
	},
}

var DefaultIndexes = map[string] /* collection */ []mongo.IndexModel{
	defaultColNameAccount:              accountIndexModels,
	defaultColNameBalance:              balanceIndexModels,
	defaultColNameOperation:            operationIndexModels,
	defaultColNameTimeStamp:            timestampIndexModels,
	defaultColNameNFTCollection:        nftCollectionIndexModels,
	defaultColNameNFT:                  nftIndexModels,
	defaultColNameNFTOperator:          nftOperatorIndexModels,
	defaultColNameDIDCredentialService: didCredentialServiceIndexModels,
	defaultColNameDIDCredential:        didCredentialIndexModels,
	defaultColNameHolder:               didCredentialHolderIndexModels,
	defaultColNameTemplate:             didCredentialTemplateIndexModels,
	defaultColNameDAO:                  daoServiceIndexModels,
	defaultColNameDAOProposal:          daoProposalIndexModels,
	defaultColNameDAODelegators:        daoDelegatorsIndexModels,
	defaultColNameDAOVoters:            daoVotersIndexModels,
	defaultColNameDAOVotingPowerBox:    daoVotingPowerBoxIndexModels,
	defaultColNameToken:                tokenServiceIndexModels,
	defaultColNameTokenBalance:         tokenBalanceIndexModels,
	defaultColNamePoint:                pointServiceIndexModels,
	defaultColNamePointBalance:         pointBalanceIndexModels,
	defaultColNameStorage:              storageDataServiceIndexModels,
	defaultColNameStorageData:          storageDataIndexModels,
	defaultColNamePrescription:         prescriptionServiceIndexModels,
	defaultColNamePrescriptionInfo:     prescriptionInfoIndexModels,
	defaultColNameDIDRegistry:          didRegistryIndexModels,
	defaultColNameDIDData:              didRegistryDataIndexModels,
	defaultColNameDIDDocument:          didRegistryDocumentIndexModels,
	defaultColNameDmile:                dmileServiceIndexModels,
	defaultColNameDmileData:            dmileServiceDataIndexModels,
}
