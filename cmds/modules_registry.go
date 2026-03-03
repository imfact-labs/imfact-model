package cmds

import (
	"github.com/imfact-labs/currency-model/app/modulekit"
	"github.com/imfact-labs/imfact-model/runtime/spec"
)

func mustBuildModuleRegistry() *modulekit.Registry {
	return spec.MustBuildModuleRegistry()
}
