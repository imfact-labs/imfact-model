package cmds

import (
	credentialcmds "github.com/ProtoconNet/mitum-credential/cmds"
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	dmilecmds "github.com/ProtoconNet/mitum-d-mile/cmds"
	daocmds "github.com/ProtoconNet/mitum-dao/cmds"
	didcmds "github.com/ProtoconNet/mitum-did-registry/cmds"
	nftcmds "github.com/ProtoconNet/mitum-nft/cmds"
	pointcmds "github.com/ProtoconNet/mitum-point/cmds"
	prescriptioncmds "github.com/ProtoconNet/mitum-prescription/cmds"
	storagecmds "github.com/ProtoconNet/mitum-storage/cmds"
	timestampcmds "github.com/ProtoconNet/mitum-timestamp/cmds"
	tokencmds "github.com/ProtoconNet/mitum-token/cmds"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

func init() {
	defaultLen := len(launch.Hinters)
	currencyExtendedLen := defaultLen + len(currencycmds.AddedHinters)
	nftExtendedLen := currencyExtendedLen + len(nftcmds.AddedHinters)
	timestampExtendedLen := nftExtendedLen + len(timestampcmds.AddedHinters)
	credentialExtendedLen := timestampExtendedLen + len(credentialcmds.AddedHinters)
	tokenExtendedLen := credentialExtendedLen + len(tokencmds.AddedHinters)
	pointExtendedLen := tokenExtendedLen + len(pointcmds.AddedHinters)
	daoExtendedLen := pointExtendedLen + len(daocmds.AddedHinters)
	storageExtendedLen := daoExtendedLen + len(storagecmds.AddedHinters)
	prescriptionExtendedLen := storageExtendedLen + len(prescriptioncmds.AddedHinters)
	didExtendedLen := prescriptionExtendedLen + len(didcmds.AddedHinters)
	allExtendedLen := didExtendedLen + len(dmilecmds.AddedHinters)

	Hinters = make([]encoder.DecodeDetail, allExtendedLen)
	copy(Hinters, launch.Hinters)
	copy(Hinters[defaultLen:currencyExtendedLen], currencycmds.AddedHinters)
	copy(Hinters[currencyExtendedLen:nftExtendedLen], nftcmds.AddedHinters)
	copy(Hinters[nftExtendedLen:timestampExtendedLen], timestampcmds.AddedHinters)
	copy(Hinters[timestampExtendedLen:credentialExtendedLen], credentialcmds.AddedHinters)
	copy(Hinters[credentialExtendedLen:tokenExtendedLen], tokencmds.AddedHinters)
	copy(Hinters[tokenExtendedLen:pointExtendedLen], pointcmds.AddedHinters)
	copy(Hinters[pointExtendedLen:daoExtendedLen], daocmds.AddedHinters)
	copy(Hinters[daoExtendedLen:storageExtendedLen], storagecmds.AddedHinters)
	copy(Hinters[storageExtendedLen:prescriptionExtendedLen], prescriptioncmds.AddedHinters)
	copy(Hinters[prescriptionExtendedLen:didExtendedLen], didcmds.AddedHinters)
	copy(Hinters[didExtendedLen:], dmilecmds.AddedHinters)

	defaultSupportedLen := len(launch.SupportedProposalOperationFactHinters)
	currencySupportedExtendedLen := defaultSupportedLen + len(currencycmds.AddedSupportedHinters)
	nftSupportedExtendedLen := currencySupportedExtendedLen + len(nftcmds.AddedSupportedHinters)
	timestampSupportedExtendedLen := nftSupportedExtendedLen + len(timestampcmds.AddedSupportedHinters)
	credentialSupportedExtendedLen := timestampSupportedExtendedLen + len(credentialcmds.AddedSupportedHinters)
	tokenSupportedExtendedLen := credentialSupportedExtendedLen + len(tokencmds.AddedSupportedHinters)
	pointSupportedExtendedLen := tokenSupportedExtendedLen + len(pointcmds.AddedSupportedHinters)
	daoSupportedExtendedLen := pointSupportedExtendedLen + len(daocmds.AddedSupportedHinters)
	storageSupportedExtendedLen := daoSupportedExtendedLen + len(storagecmds.AddedSupportedHinters)
	prescriptionSupportedExtendedLen := storageSupportedExtendedLen + len(prescriptioncmds.AddedSupportedHinters)
	didSupportedExtendedLen := prescriptionSupportedExtendedLen + len(didcmds.AddedSupportedHinters)
	allSupportedExtendedLen := didSupportedExtendedLen + len(dmilecmds.AddedSupportedHinters)

	SupportedProposalOperationFactHinters = make(
		[]encoder.DecodeDetail,
		allSupportedExtendedLen)
	copy(SupportedProposalOperationFactHinters, launch.SupportedProposalOperationFactHinters)
	copy(SupportedProposalOperationFactHinters[defaultSupportedLen:currencySupportedExtendedLen], currencycmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[currencySupportedExtendedLen:nftSupportedExtendedLen], nftcmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[nftSupportedExtendedLen:timestampSupportedExtendedLen], timestampcmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[timestampSupportedExtendedLen:credentialSupportedExtendedLen], credentialcmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[credentialSupportedExtendedLen:tokenSupportedExtendedLen], tokencmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[tokenSupportedExtendedLen:pointSupportedExtendedLen], pointcmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[pointSupportedExtendedLen:daoSupportedExtendedLen], daocmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[daoSupportedExtendedLen:storageSupportedExtendedLen], storagecmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[storageSupportedExtendedLen:prescriptionSupportedExtendedLen], prescriptioncmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[prescriptionSupportedExtendedLen:didSupportedExtendedLen], didcmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[didSupportedExtendedLen:], dmilecmds.AddedSupportedHinters)
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
