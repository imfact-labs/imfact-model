package digest

import (
	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	utilc "github.com/ProtoconNet/mitum-currency/v3/digest/util"
	"github.com/ProtoconNet/mitum-prescription/state"
	"github.com/ProtoconNet/mitum-prescription/types"
	"github.com/ProtoconNet/mitum2/base"
	utilm "github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func PrescriptionDesign(st *cdigest.Database, contract string) (types.Design, base.State, error) {
	filter := utilc.NewBSONFilter("contract", contract)
	q := filter.D()

	opt := options.FindOne().SetSort(
		utilc.NewBSONFilter("height", -1).D(),
	)
	var sta base.State
	if err := st.MongoClient().GetByFilter(
		defaultColNamePrescription,
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
		return types.Design{}, nil, utilm.ErrNotFound.WithMessage(err, "prescription design by contract account %v", contract)
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

func PrescriptionInfo(db *cdigest.Database, contract, hash string) (*types.PrescriptionInfo, error) {
	filter := utilc.NewBSONFilter("contract", contract)
	filter = filter.Add("prescription_hash", hash)
	q := filter.D()

	opt := options.FindOne().SetSort(
		utilc.NewBSONFilter("height", -1).D(),
	)
	var data *types.PrescriptionInfo
	var sta base.State
	var err error
	if err := db.MongoClient().GetByFilter(
		defaultColNamePrescriptionInfo,
		q,
		func(res *mongo.SingleResult) error {
			sta, err = cdigest.LoadState(res.Decode, db.Encoders())
			if err != nil {
				return err
			}
			i, err := state.GetPrescriptionInfoFromState(sta)
			if err != nil {
				return err
			}
			data = &i
			return nil
		},
		opt,
	); err != nil {
		return nil, utilm.ErrNotFound.WithMessage(err, "prescription info for key %s in contract account %s", hash, contract)
	}

	if data != nil {
		return data, nil
	} else {
		return nil, errors.Errorf("data is nil")
	}
}
