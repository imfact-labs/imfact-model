package steps

import (
	"context"

	"github.com/imfact-labs/imfact-model/runtime/contracts"
	"github.com/imfact-labs/imfact-model/runtime/spec"
	"github.com/imfact-labs/mitum2/launch"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/imfact-labs/mitum2/util/hint"
	"github.com/pkg/errors"
)

func PAddHinters(pctx context.Context) (context.Context, error) {
	e := util.StringError("add hinters")

	var encs *encoder.Encoders
	var f contracts.ProposalOperationFactHintFunc = IsSupportedProposalOperationFactHintFunc

	if err := util.LoadFromContextOK(pctx, launch.EncodersContextKey, &encs); err != nil {
		return pctx, e.Wrap(err)
	}
	pctx = context.WithValue(pctx, contracts.ProposalOperationFactHintContextKey, f)

	if err := LoadHinters(encs); err != nil {
		return pctx, e.Wrap(err)
	}

	return pctx, nil
}

func LoadHinters(encs *encoder.Encoders) error {
	for i := range spec.Hinters {
		if err := encs.AddDetail(spec.Hinters[i]); err != nil {
			return errors.Wrap(err, "add hinter to encoder")
		}
	}

	for i := range spec.SupportedProposalOperationFactHinters {
		if err := encs.AddDetail(spec.SupportedProposalOperationFactHinters[i]); err != nil {
			return errors.Wrap(err, "add supported proposal operation fact hinter to encoder")
		}
	}

	return nil
}

func IsSupportedProposalOperationFactHintFunc() func(hint.Hint) bool {
	return func(ht hint.Hint) bool {
		for i := range spec.SupportedProposalOperationFactHinters {
			s := spec.SupportedProposalOperationFactHinters[i].Hint
			if ht.Type() != s.Type() {
				continue
			}

			return ht.IsCompatible(s)
		}

		return false
	}
}
