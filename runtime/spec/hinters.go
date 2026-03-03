package spec

import "github.com/imfact-labs/mitum2/util/encoder"

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

func init() {
	registry := MustBuildModuleRegistry()
	entries := registry.Entries()

	for i := range entries {
		Hinters = append(Hinters, entries[i].Hinters...)
		SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, entries[i].SupportedFacts...)
	}
}
