package spec

import (
	"sync"

	ccmodule "github.com/imfact-labs/currency-model/app/module"
	"github.com/imfact-labs/currency-model/app/modulekit"
	daomodule "github.com/imfact-labs/dao-model/module"
	nmodule "github.com/imfact-labs/nft-model/module"
	pmmodule "github.com/imfact-labs/payment-model/module"
	smodule "github.com/imfact-labs/storage-model/module"
	tsmodule "github.com/imfact-labs/timestamp-model/module"
	tkmodule "github.com/imfact-labs/token-model/module"
)

var composedModules = []modulekit.ModelModule{
	ccmodule.Module{},
	nmodule.Module{},
	tsmodule.Module{},
	tkmodule.Module{},
	daomodule.Module{},
	smodule.Module{},
	pmmodule.Module{},
}

var (
	moduleRegistryOnce sync.Once
	moduleRegistry     *modulekit.Registry
	moduleRegistryErr  error
)

func LoadModuleRegistry() (*modulekit.Registry, error) {
	moduleRegistryOnce.Do(func() {
		moduleRegistry, moduleRegistryErr = buildModuleRegistry()
	})

	return moduleRegistry, moduleRegistryErr
}

func buildModuleRegistry() (*modulekit.Registry, error) {
	registry := modulekit.NewRegistry()

	for i := range composedModules {
		if err := registry.Register(composedModules[i]); err != nil {
			return nil, err
		}
	}

	for i := range composedModules {
		if err := registry.ValidateModuleContract(composedModules[i].ID()); err != nil {
			return nil, err
		}
	}

	return registry, nil
}

func MustBuildModuleRegistry() *modulekit.Registry {
	registry, err := LoadModuleRegistry()
	if err != nil {
		panic(err)
	}

	return registry
}
