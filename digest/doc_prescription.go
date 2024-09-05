package digest

import (
	mongodb "github.com/ProtoconNet/mitum-currency/v3/digest/mongodb"
	bsonutil "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	cstate "github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum-prescription/state"
	"github.com/ProtoconNet/mitum-prescription/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type PrescriptionDesignDoc struct {
	mongodb.BaseDoc
	st     base.State
	design types.Design
}

// NewPrescriptionDesignDoc get the state of prescription Design
func NewPrescriptionDesignDoc(st base.State, enc encoder.Encoder) (PrescriptionDesignDoc, error) {
	design, err := state.GetDesignFromState(st)

	if err != nil {
		return PrescriptionDesignDoc{}, err
	}

	b, err := mongodb.NewBaseDoc(nil, st, enc)
	if err != nil {
		return PrescriptionDesignDoc{}, err
	}

	return PrescriptionDesignDoc{
		BaseDoc: b,
		st:      st,
		design:  design,
	}, nil
}

func (doc PrescriptionDesignDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := cstate.ParseStateKey(doc.st.Key(), state.PrescriptionStateKeyPrefix, 3)

	m["contract"] = parsedKey[1]
	m["height"] = doc.st.Height()

	return bsonutil.Marshal(m)
}

type PrescriptionInfoDoc struct {
	mongodb.BaseDoc
	st   base.State
	info types.PrescriptionInfo
}

func NewPrescriptionInfoDoc(st base.State, enc encoder.Encoder) (PrescriptionInfoDoc, error) {
	info, err := state.GetPrescriptionInfoFromState(st)
	if err != nil {
		return PrescriptionInfoDoc{}, err
	}

	b, err := mongodb.NewBaseDoc(nil, st, enc)
	if err != nil {
		return PrescriptionInfoDoc{}, err
	}

	return PrescriptionInfoDoc{
		BaseDoc: b,
		st:      st,
		info:    info,
	}, nil
}

func (doc PrescriptionInfoDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := cstate.ParseStateKey(doc.st.Key(), state.PrescriptionStateKeyPrefix, 4)
	if err != nil {
		return nil, err
	}

	m["contract"] = parsedKey[1]
	m["prescription_hash"] = doc.info.PrescriptionHash()
	m["height"] = doc.st.Height()

	return bsonutil.Marshal(m)
}
