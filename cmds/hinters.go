package cmds

import (
	crcmds "github.com/ProtoconNet/mitum-credential/cmds"
	ccmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	daocmds "github.com/ProtoconNet/mitum-dao/cmds"
	ncmds "github.com/ProtoconNet/mitum-nft/cmds"
	pmcmds "github.com/ProtoconNet/mitum-payment/cmds"
	pcmds "github.com/ProtoconNet/mitum-point/cmds"
	scmds "github.com/ProtoconNet/mitum-storage/cmds"
	tscmds "github.com/ProtoconNet/mitum-timestamp/cmds"
	tkcmds "github.com/ProtoconNet/mitum-token/cmds"
	"github.com/ProtoconNet/mitum2/util/encoder"
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
	Hinters = append(Hinters, pcmds.AddedHinters...)
	Hinters = append(Hinters, daocmds.AddedHinters...)
	Hinters = append(Hinters, scmds.AddedHinters...)
	Hinters = append(Hinters, pmcmds.AddedHinters...)

	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, ccmds.SupportedProposalOperationFactHinters...)
	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, ncmds.AddedSupportedHinters...)
	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, tscmds.AddedSupportedHinters...)
	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, crcmds.AddedSupportedHinters...)
	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, tkcmds.AddedSupportedHinters...)
	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, pcmds.AddedSupportedHinters...)
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
