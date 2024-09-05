package digest

import (
	"github.com/ProtoconNet/mitum-prescription/state"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/mongo"
)

func (bs *BlockSession) preparePrescription() error {
	if len(bs.sts) < 1 {
		return nil
	}

	var prescriptionModels []mongo.WriteModel
	var prescriptionInfoModels []mongo.WriteModel
	for i := range bs.sts {
		st := bs.sts[i]
		switch {
		case state.IsDesignStateKey(st.Key()):
			j, err := bs.handlePrescriptionDesignState(st)
			if err != nil {
				return err
			}
			prescriptionModels = append(prescriptionModels, j...)
		case state.IsPrescriptionInfoStateKey(st.Key()):
			j, err := bs.handlePrescriptionInfoState(st)
			if err != nil {
				return err
			}
			prescriptionInfoModels = append(prescriptionInfoModels, j...)
		default:
			continue
		}
	}

	bs.prescriptionModels = prescriptionModels
	bs.prescriptionInfoDataModels = prescriptionInfoModels

	return nil
}

func (bs *BlockSession) handlePrescriptionDesignState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if prescriptionDesignDoc, err := NewPrescriptionDesignDoc(st, bs.st.Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(prescriptionDesignDoc),
		}, nil
	}
}

func (bs *BlockSession) handlePrescriptionInfoState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if prescriptionInfoDoc, err := NewPrescriptionInfoDoc(st, bs.st.Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(prescriptionInfoDoc),
		}, nil
	}
}
