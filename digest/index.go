package digest

import (
	crdigest "github.com/ProtoconNet/mitum-credential/digest"
	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	daodigest "github.com/ProtoconNet/mitum-dao/digest"
	ndigest "github.com/ProtoconNet/mitum-nft/digest"
	pmdigest "github.com/ProtoconNet/mitum-payment/digest"
	pdigest "github.com/ProtoconNet/mitum-point/digest"
	sdigest "github.com/ProtoconNet/mitum-storage/digest"
	tsdigest "github.com/ProtoconNet/mitum-timestamp/digest"
	tkdigest "github.com/ProtoconNet/mitum-token/digest"
	"go.mongodb.org/mongo-driver/mongo"
)

var AllIndexes = []map[string][]mongo.IndexModel{
	cdigest.DefaultIndexes,
	crdigest.DefaultIndexes,
	tsdigest.DefaultIndexes,
	ndigest.DefaultIndexes,
	sdigest.DefaultIndexes,
	tkdigest.DefaultIndexes,
	pdigest.DefaultIndexes,
	daodigest.DefaultIndexes,
	pmdigest.DefaultIndexes,
}

var DefaultIndexes = cdigest.DefaultIndexes

func init() {
	for i := range AllIndexes {
		for k, v := range AllIndexes[i] {
			if _, found := DefaultIndexes[k]; !found {
				DefaultIndexes[k] = v
			}
		}
	}
}
