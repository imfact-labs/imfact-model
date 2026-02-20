package cmds

import (
	"context"

	ccmds "github.com/imfact-labs/currency-model/app/cmds"
	cprocessor "github.com/imfact-labs/currency-model/operation/processor"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/ps"
)

var PNameOperationProcessorsMap = ps.Name("mitum-minic-operation-processors-map")

func POperationProcessorsMap(pctx context.Context) (context.Context, error) {
	var opr *cprocessor.OperationProcessor

	if err := util.LoadFromContextOK(pctx,
		ccmds.OperationProcessorContextKey, &opr,
	); err != nil {
		return pctx, err
	}

	err := opr.SetGetNewProcessorFunc(cprocessor.GetNewProcessor)
	if err != nil {
		return pctx, err
	}

	return pctx, nil
}
