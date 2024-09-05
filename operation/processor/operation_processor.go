package processor

import (
	"fmt"
	"github.com/ProtoconNet/mitum-credential/operation/credential"
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/operation/extension"
	currencyprocessor "github.com/ProtoconNet/mitum-currency/v3/operation/processor"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/operation/nft"
	"github.com/ProtoconNet/mitum-point/operation/point"
	"github.com/ProtoconNet/mitum-prescription/operation/prescription"
	"github.com/ProtoconNet/mitum-storage/operation/storage"
	"github.com/ProtoconNet/mitum-timestamp/operation/timestamp"
	"github.com/ProtoconNet/mitum-token/operation/token"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

const (
	DuplicationTypeSender           currencytypes.DuplicationType = "sender"
	DuplicationTypeCurrency         currencytypes.DuplicationType = "currency"
	DuplicationTypeContract         currencytypes.DuplicationType = "contract"
	DuplicationTypeCredential       currencytypes.DuplicationType = "credential"
	DuplicationTypeStorageData      currencytypes.DuplicationType = "storagedata"
	DuplicationTypePrescriptionInfo currencytypes.DuplicationType = "prescriptioninfo"
)

func CheckDuplication(opr *currencyprocessor.OperationProcessor, op base.Operation) error {
	opr.Lock()
	defer opr.Unlock()

	var duplicationTypeSenderID string
	var duplicationTypeCurrencyID string
	var duplicationTypeCredentialID []string
	var duplicationTypeContractID string
	var duplicationTypeStorageData string
	var duplicationTypePrescriptionInfo string
	var newAddresses []base.Address

	switch t := op.(type) {
	case currency.CreateAccount:
		fact, ok := t.Fact().(currency.CreateAccountFact)
		if !ok {
			return errors.Errorf("expected CreateAccountFact, not %T", t.Fact())
		}
		as, err := fact.Targets()
		if err != nil {
			return errors.Errorf("failed to get Addresses")
		}
		newAddresses = as
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case currency.UpdateKey:
		fact, ok := t.Fact().(currency.UpdateKeyFact)
		if !ok {
			return errors.Errorf("expected UpdateKeyFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case currency.Transfer:
		fact, ok := t.Fact().(currency.TransferFact)
		if !ok {
			return errors.Errorf("expected TransferFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case currency.RegisterCurrency:
		fact, ok := t.Fact().(currency.RegisterCurrencyFact)
		if !ok {
			return errors.Errorf("expected RegisterCurrencyFact, not %T", t.Fact())
		}
		duplicationTypeCurrencyID = currencyprocessor.DuplicationKey(fact.Currency().Currency().String(), DuplicationTypeCurrency)
	case currency.UpdateCurrency:
		fact, ok := t.Fact().(currency.UpdateCurrencyFact)
		if !ok {
			return errors.Errorf("expected UpdateCurrencyFact, not %T", t.Fact())
		}
		duplicationTypeCurrencyID = currencyprocessor.DuplicationKey(fact.Currency().String(), DuplicationTypeCurrency)
	case currency.Mint:
	case extensioncurrency.CreateContractAccount:
		fact, ok := t.Fact().(extensioncurrency.CreateContractAccountFact)
		if !ok {
			return errors.Errorf("expected CreateContractAccountFact, not %T", t.Fact())
		}
		as, err := fact.Targets()
		if err != nil {
			return errors.Errorf("failed to get Addresses")
		}
		newAddresses = as
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
		duplicationTypeContractID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeContract)
	case extensioncurrency.Withdraw:
		fact, ok := t.Fact().(extensioncurrency.WithdrawFact)
		if !ok {
			return errors.Errorf("expected WithdrawFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case nft.RegisterModel:
		fact, ok := t.Fact().(nft.RegisterModelFact)
		if !ok {
			return errors.Errorf("expected RegisterModelFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
		duplicationTypeContractID = currencyprocessor.DuplicationKey(fact.Contract().String(), DuplicationTypeContract)
	case nft.UpdateModelConfig:
		fact, ok := t.Fact().(nft.UpdateModelConfigFact)
		if !ok {
			return errors.Errorf("expected UpdateModelConfigFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case nft.Mint:
		fact, ok := t.Fact().(nft.MintFact)
		if !ok {
			return errors.Errorf("expected MintFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case nft.Transfer:
		fact, ok := t.Fact().(nft.TransferFact)
		if !ok {
			return errors.Errorf("expected TransferFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case nft.ApproveAll:
		fact, ok := t.Fact().(nft.ApproveAllFact)
		if !ok {
			return errors.Errorf("expected ApproveAllFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case nft.Approve:
		fact, ok := t.Fact().(nft.ApproveFact)
		if !ok {
			return errors.Errorf("expected ApproveFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case nft.AddSignature:
		fact, ok := t.Fact().(nft.AddSignatureFact)
		if !ok {
			return errors.Errorf("expected AddSignatureFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case timestamp.RegisterModel:
		fact, ok := t.Fact().(timestamp.RegisterModelFact)
		if !ok {
			return errors.Errorf("expected RegisterModelFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
		duplicationTypeContractID = currencyprocessor.DuplicationKey(fact.Contract().String(), DuplicationTypeContract)
	case timestamp.Issue:
		fact, ok := t.Fact().(timestamp.IssueFact)
		if !ok {
			return errors.Errorf("expected IssueFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case credential.RegisterModel:
		fact, ok := t.Fact().(credential.RegisterModelFact)
		if !ok {
			return errors.Errorf("expected RegisterModelFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
		duplicationTypeContractID = currencyprocessor.DuplicationKey(fact.Contract().String(), DuplicationTypeContract)
	case credential.AddTemplate:
		fact, ok := t.Fact().(credential.AddTemplateFact)
		if !ok {
			return errors.Errorf("expected AddTemplateFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case credential.Issue:
		fact, ok := t.Fact().(credential.IssueFact)
		if !ok {
			return errors.Errorf("expected IssueFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
		var credentials []string
		for _, v := range fact.Items() {
			key := currencyprocessor.DuplicationKey(fmt.Sprintf("%s:%s:%s", v.Contract().String(), v.TemplateID(), v.CredentialID()), DuplicationTypeCredential)
			credentials = append(credentials, key)
		}
		duplicationTypeCredentialID = credentials
	case credential.Revoke:
		fact, ok := t.Fact().(credential.RevokeFact)
		if !ok {
			return errors.Errorf("expected RevokeFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
		var credentials []string
		for _, v := range fact.Items() {
			key := currencyprocessor.DuplicationKey(fmt.Sprintf("%s:%s:%s", v.Contract().String(), v.TemplateID(), v.CredentialID()), DuplicationTypeCredential)
			credentials = append(credentials, key)
		}
		duplicationTypeCredentialID = credentials
	case token.Mint:
		fact, ok := t.Fact().(token.MintFact)
		if !ok {
			return errors.Errorf("expected MintFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case token.RegisterModel:
		fact, ok := t.Fact().(token.RegisterModelFact)
		if !ok {
			return errors.Errorf("expected RegisterModelFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
		duplicationTypeContractID = currencyprocessor.DuplicationKey(fact.Contract().String(), DuplicationTypeContract)
	case token.Burn:
		fact, ok := t.Fact().(token.BurnFact)
		if !ok {
			return errors.Errorf("expected BurnFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case token.Approve:
		fact, ok := t.Fact().(token.ApproveFact)
		if !ok {
			return errors.Errorf("expected ApproveFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case token.Transfer:
		fact, ok := t.Fact().(token.TransferFact)
		if !ok {
			return errors.Errorf("expected TransferFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case token.TransferFrom:
		fact, ok := t.Fact().(token.TransferFromFact)
		if !ok {
			return errors.Errorf("expected TransferFromFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case point.Mint:
		fact, ok := t.Fact().(point.MintFact)
		if !ok {
			return errors.Errorf("expected MintFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case point.RegisterModel:
		fact, ok := t.Fact().(point.RegisterModelFact)
		if !ok {
			return errors.Errorf("expected RegisterModelFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
		duplicationTypeContractID = currencyprocessor.DuplicationKey(fact.Contract().String(), DuplicationTypeContract)
	case point.Burn:
		fact, ok := t.Fact().(point.BurnFact)
		if !ok {
			return errors.Errorf("expected BurnFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case point.Approve:
		fact, ok := t.Fact().(point.ApproveFact)
		if !ok {
			return errors.Errorf("expected ApproveFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case point.Transfer:
		fact, ok := t.Fact().(point.TransferFact)
		if !ok {
			return errors.Errorf("expected TransferFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case point.TransferFrom:
		fact, ok := t.Fact().(point.TransferFromFact)
		if !ok {
			return errors.Errorf("expected TransferFromFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
	case storage.RegisterModel:
		fact, ok := t.Fact().(storage.RegisterModelFact)
		if !ok {
			return errors.Errorf("expected RegisterModelFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
		duplicationTypeContractID = currencyprocessor.DuplicationKey(fact.Contract().String(), DuplicationTypeContract)
	case storage.CreateData:
		fact, ok := t.Fact().(storage.CreateDataFact)
		if !ok {
			return errors.Errorf("expected CreateDataFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
		duplicationTypeStorageData = currencyprocessor.DuplicationKey(
			fmt.Sprintf("%s:%s", fact.Contract().String(), fact.DataKey()), DuplicationTypeStorageData)
	case storage.UpdateData:
		fact, ok := t.Fact().(storage.UpdateDataFact)
		if !ok {
			return errors.Errorf("expected UpdateDataFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
		duplicationTypeStorageData = currencyprocessor.DuplicationKey(
			fmt.Sprintf("%s:%s", fact.Contract().String(), fact.DataKey()), DuplicationTypeStorageData)
	case storage.DeleteData:
		fact, ok := t.Fact().(storage.DeleteDataFact)
		if !ok {
			return errors.Errorf("expected DeleteDataFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
		duplicationTypeStorageData = currencyprocessor.DuplicationKey(
			fmt.Sprintf("%s:%s", fact.Contract().String(), fact.DataKey()), DuplicationTypeStorageData)
	case prescription.RegisterModel:
		fact, ok := t.Fact().(prescription.RegisterModelFact)
		if !ok {
			return errors.Errorf("expected RegisterModelFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
		duplicationTypeContractID = currencyprocessor.DuplicationKey(fact.Contract().String(), DuplicationTypeContract)
	case prescription.RegisterPrescription:
		fact, ok := t.Fact().(prescription.RegisterPrescriptionFact)
		if !ok {
			return errors.Errorf("expected RegisterPrescriptionFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
		duplicationTypePrescriptionInfo = currencyprocessor.DuplicationKey(
			fmt.Sprintf("%s:%s", fact.Contract().String(), fact.PrescriptionHash()), DuplicationTypePrescriptionInfo)
	case prescription.UsePrescription:
		fact, ok := t.Fact().(prescription.UsePrescriptionFact)
		if !ok {
			return errors.Errorf("expected UsePrescriptionFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = currencyprocessor.DuplicationKey(fact.Sender().String(), DuplicationTypeSender)
		duplicationTypePrescriptionInfo = currencyprocessor.DuplicationKey(
			fmt.Sprintf("%s:%s", fact.Contract().String(), fact.PrescriptionHash()), DuplicationTypePrescriptionInfo)
	default:
		return nil
	}

	if len(duplicationTypeSenderID) > 0 {
		if _, found := opr.Duplicated[duplicationTypeSenderID]; found {
			return errors.Errorf("proposal cannot have duplicated sender, %v", duplicationTypeSenderID)
		}

		opr.Duplicated[duplicationTypeSenderID] = struct{}{}
	}

	if len(duplicationTypeCurrencyID) > 0 {
		if _, found := opr.Duplicated[duplicationTypeCurrencyID]; found {
			return errors.Errorf(
				"cannot register duplicated currency id, %v within a proposal",
				duplicationTypeCurrencyID,
			)
		}

		opr.Duplicated[duplicationTypeCurrencyID] = struct{}{}
	}
	if len(duplicationTypeContractID) > 0 {
		if _, found := opr.Duplicated[duplicationTypeContractID]; found {
			return errors.Errorf(
				"cannot use a duplicated contract for registering in contract model, %v within a proposal",
				duplicationTypeSenderID,
			)
		}

		opr.Duplicated[duplicationTypeContractID] = struct{}{}
	}
	if len(duplicationTypeCredentialID) > 0 {
		for _, v := range duplicationTypeCredentialID {
			if _, found := opr.Duplicated[v]; found {
				return errors.Errorf(
					"cannot use a duplicated contract-template-credential for credential, %v within a proposal",
					v,
				)
			}
			opr.Duplicated[v] = struct{}{}
		}
	}

	if len(duplicationTypeStorageData) > 0 {
		if _, found := opr.Duplicated[duplicationTypeStorageData]; found {
			return errors.Errorf(
				"cannot use a duplicated contract-key for storage data, %v within a proposal",
				duplicationTypeStorageData,
			)
		}

		opr.Duplicated[duplicationTypeStorageData] = struct{}{}
	}

	if len(duplicationTypePrescriptionInfo) > 0 {
		if _, found := opr.Duplicated[duplicationTypePrescriptionInfo]; found {
			return errors.Errorf(
				"cannot use a duplicated contract-hash for prescription info, %v within a proposal",
				duplicationTypeStorageData,
			)
		}

		opr.Duplicated[duplicationTypeStorageData] = struct{}{}
	}

	if len(newAddresses) > 0 {
		if err := opr.CheckNewAddressDuplication(newAddresses); err != nil {
			return err
		}
	}

	return nil
}
