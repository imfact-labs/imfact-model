package cmds

import (
	crcmds "github.com/imfact-labs/credential-model/cmds"
	ccmds "github.com/imfact-labs/currency-model/app/cmds"
	daocmds "github.com/imfact-labs/dao-model/cmds"
	"github.com/imfact-labs/mitum2/util/encoder"
	ncmds "github.com/imfact-labs/nft-model/cmds"
	pmcmds "github.com/imfact-labs/payment-model/cmds"
	scmds "github.com/imfact-labs/storage-model/cmds"
	tscmds "github.com/imfact-labs/timestamp-model/cmds"
	tkcmds "github.com/imfact-labs/token-model/cmds"
	"github.com/pkg/errors"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

func init() {
	Hinters = append(Hinters, ccmds.Hinters...)
	Hinters = append(Hinters, ncmds.AddedHinters...)
	Hinters = append(Hinters, tscmds.AddedHinters...)
	Hinters = append(Hinters, crcmds.AddedHinters...)
	Hinters = append(Hinters, tkcmds.AddedHinters...)
	Hinters = append(Hinters, daocmds.AddedHinters...)
	Hinters = append(Hinters, scmds.AddedHinters...)
	Hinters = append(Hinters, pmcmds.AddedHinters...)

	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, ccmds.SupportedProposalOperationFactHinters...)
	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, ncmds.AddedSupportedHinters...)
	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, tscmds.AddedSupportedHinters...)
	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, crcmds.AddedSupportedHinters...)
	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, tkcmds.AddedSupportedHinters...)
	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, daocmds.AddedSupportedHinters...)
	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, scmds.AddedSupportedHinters...)
	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, pmcmds.AddedSupportedHinters...)
}

func LoadHinters(encs *encoder.Encoders) error {
	for i := range Hinters {
		if err := encs.AddDetail(Hinters[i]); err != nil {
			return errors.Wrap(err, "add hinter to encoder")
		}
	}

	for i := range SupportedProposalOperationFactHinters {
		if err := encs.AddDetail(SupportedProposalOperationFactHinters[i]); err != nil {
			return errors.Wrap(err, "add supported proposal operation fact hinter to encoder")
		}
	}

	return nil
}
