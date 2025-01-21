package cmds

import (
	"context"

	ccmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	cprocessor "github.com/ProtoconNet/mitum-currency/v3/operation/processor"
	"github.com/ProtoconNet/mitum-minic/operation/processor"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/ps"
)

var PNameOperationProcessorsMap = ps.Name("mitum-minic-operation-processors-map")

func POperationProcessorsMap(pctx context.Context) (context.Context, error) {
	var opr *cprocessor.OperationProcessor

	if err := util.LoadFromContextOK(pctx,
		ccmds.OperationProcessorContextKey, &opr,
	); err != nil {
		return pctx, err
	}

	err := opr.SetCheckDuplicationFunc(processor.CheckDuplication)
	if err != nil {
		return pctx, err
	}

	return pctx, nil
}
