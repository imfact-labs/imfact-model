package cmds

import (
	ccmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util/ps"
)

var (
	PNameDigest           = ps.Name("digest")
	PNameDigestStart      = ps.Name("digest_star")
	PNameMongoDBsDataBase = ps.Name("mongodb_database")
)

func DefaultRunPS() *ps.PS {
	pps := ps.NewPS("cmd-run")

	_ = pps.
		AddOK(launch.PNameEncoder, ccmds.PEncoder, nil).
		AddOK(launch.PNameDesign, launch.PLoadDesign, nil, launch.PNameEncoder).
		AddOK(ccmds.PNameDigestDesign, ccmds.PLoadDigestDesign, nil, launch.PNameEncoder).
		AddOK(launch.PNameTimeSyncer, launch.PStartTimeSyncer, launch.PCloseTimeSyncer, launch.PNameDesign).
		AddOK(launch.PNameLocal, launch.PLocal, nil, launch.PNameDesign).
		AddOK(launch.PNameBlockItemReaders, launch.PBlockItemReaders, nil, launch.PNameDesign).
		AddOK(launch.PNameStorage, launch.PStorage, nil, launch.PNameLocal, launch.PNameBlockItemReaders).
		AddOK(launch.PNameProposalMaker, ccmds.PProposalMaker, nil, launch.PNameStorage).
		AddOK(launch.PNameNetwork, launch.PNetwork, nil, launch.PNameStorage).
		AddOK(launch.PNameMemberlist, ccmds.PMemberlist, nil, launch.PNameNetwork).
		AddOK(launch.PNameStartStorage, launch.PStartStorage, launch.PCloseStorage, launch.PNameStartNetwork).
		AddOK(launch.PNameStartNetwork, launch.PStartNetwork, launch.PCloseNetwork, launch.PNameStates).
		AddOK(launch.PNameStartMemberlist, ccmds.PStartMemberlist, ccmds.PCloseMemberlist, launch.PNameStartNetwork).
		AddOK(launch.PNameStartSyncSourceChecker, launch.PStartSyncSourceChecker, launch.PCloseSyncSourceChecker, launch.PNameStartNetwork).
		AddOK(launch.PNameStartLastConsensusNodesWatcher,
			launch.PStartLastConsensusNodesWatcher, launch.PCloseLastConsensusNodesWatcher, launch.PNameStartNetwork).
		AddOK(launch.PNameStates, launch.PStates, nil, launch.PNameNetwork).
		AddOK(launch.PNameStatesReady, nil, launch.PCloseStates,
			launch.PNameStartStorage,
			launch.PNameStartSyncSourceChecker,
			launch.PNameStartLastConsensusNodesWatcher,
			launch.PNameStartMemberlist,
			launch.PNameStartNetwork,
			launch.PNameStates).
		AddOK(ccmds.PNameDigesterDataBase, ccmds.ProcessDigesterDatabase, nil, ccmds.PNameDigestDesign, launch.PNameStorage).
		AddOK(ccmds.PNameAPI, ccmds.ProcessAPI, nil, ccmds.PNameDigestDesign, ccmds.PNameDigesterDataBase, launch.PNameMemberlist).
		AddOK(ccmds.PNameStartAPI, ccmds.ProcessStartAPI, nil, ccmds.PNameAPI)

	_ = pps.POK(launch.PNameDesign).
		PostAddOK(launch.PNameCheckDesign, launch.PCheckDesign).
		PostAddOK(launch.PNameINITObjectCache, launch.PINITObjectCache)

	_ = pps.POK(launch.PNameLocal).
		PostAddOK(launch.PNameDiscoveryFlag, launch.PDiscoveryFlag).
		PostAddOK(launch.PNameLoadACL, launch.PLoadACL)

	_ = pps.POK(launch.PNameBlockItemReaders).
		PreAddOK(launch.PNameBlockItemReadersDecompressFunc, launch.PBlockItemReadersDecompressFunc).
		PostAddOK(launch.PNameRemotesBlockItemReaderFunc, launch.PRemotesBlockItemReaderFunc)

	_ = pps.POK(launch.PNameStorage).
		PreAddOK(launch.PNameCheckLocalFS, launch.PCheckAndCreateLocalFS).
		PreAddOK(launch.PNameLoadDatabase, launch.PLoadDatabase).
		PostAddOK(launch.PNameCheckLeveldbStorage, launch.PCheckLeveldbStorage).
		PostAddOK(launch.PNameLoadFromDatabase, launch.PLoadFromDatabase).
		PostAddOK(launch.PNameCheckBlocksOfStorage, launch.PCheckBlocksOfStorage).
		PostAddOK(launch.PNamePatchBlockItemReaders, launch.PPatchBlockItemReaders).
		PostAddOK(launch.PNameNodeInfo, launch.PNodeInfo).
		PostAddOK(launch.PNameNodeMetric, launch.PNodeMetric)

	_ = pps.POK(launch.PNameNetwork).
		PreAddOK(launch.PNameQuicstreamClient, launch.PQuicstreamClient).
		PostAddOK(launch.PNameSyncSourceChecker, launch.PSyncSourceChecker).
		PostAddOK(launch.PNameSuffrageCandidateLimiterSet, ccmds.PSuffrageCandidateLimiterSet)

	_ = pps.POK(launch.PNameMemberlist).
		PreAddOK(launch.PNameLastConsensusNodesWatcher, launch.PLastConsensusNodesWatcher).
		PreAddOK(launch.PNameRateLimiterContextKey, launch.PNetworkRateLimiter).
		PostAddOK(launch.PNameBallotbox, launch.PBallotbox).
		PostAddOK(launch.PNameLongRunningMemberlistJoin, ccmds.PLongRunningMemberlistJoin).
		PostAddOK(launch.PNameSuffrageVoting, launch.PSuffrageVoting).
		PostAddOK(launch.PNameEventLoggingNetworkHandlers, launch.PEventLoggingNetworkHandlers)

	_ = pps.POK(launch.PNameStates).
		PreAddOK(launch.PNameProposerSelector, launch.PProposerSelector).
		PreAddOK(launch.PNameOperationProcessorsMap, ccmds.POperationProcessorsMap).
		PreAddOK(launch.PNameNetworkHandlers, ccmds.PNetworkHandlers).
		PreAddOK(launch.PNameNodeInConsensusNodesFunc, launch.PNodeInConsensusNodesFunc).
		PreAddOK(launch.PNameProposalProcessors, ccmds.PProposalProcessors).
		PreAddOK(launch.PNameBallotStuckResolver, launch.PBallotStuckResolver).
		PostAddOK(launch.PNamePatchLastConsensusNodesWatcher, launch.PPatchLastConsensusNodesWatcher).
		PostAddOK(launch.PNameStatesSetHandlers, launch.PStatesSetHandlers).
		PostAddOK(launch.PNameNetworkHandlersReadWriteNode, ccmds.PNetworkHandlersReadWriteNode).
		PostAddOK(launch.PNamePatchMemberlist, ccmds.PPatchMemberlist).
		PostAddOK(launch.PNameStatesNetworkHandlers, ccmds.PStatesNetworkHandlers).
		PostAddOK(launch.PNameHandoverNetworkHandlers, launch.PHandoverNetworkHandlers)

	return pps
}
