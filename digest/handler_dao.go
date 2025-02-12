package digest

import (
	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	daodigest "github.com/ProtoconNet/mitum-dao/digest"
	"github.com/ProtoconNet/mitum-dao/state"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/util"
	"net/http"
)

func (hd *Handlers) handleDAOService(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cacheKey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleDAODesignInGroup(contract)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			cdigest.HTTP2WriteCache(w, cacheKey, hd.expireShortLived)
		}
	}
}

func (hd *Handlers) handleDAODesignInGroup(contract string) (interface{}, error) {
	switch design, err := daodigest.DAOService(hd.database, contract); {
	case err != nil:
		return nil, util.ErrNotFound.WithMessage(err, "dao service, contract %s", contract)
	case design == nil:
		return nil, util.ErrNotFound.Errorf("dao service, contract %s", contract)
	default:
		hal, err := hd.buildDAODesignHal(contract, *design)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildDAODesignHal(contract string, design types.Design) (cdigest.Hal, error) {
	h, err := hd.combineURL(daodigest.HandlerPathDAOService, "contract", contract)
	if err != nil {
		return nil, err
	}

	hal := cdigest.NewBaseHal(design, cdigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleDAOProposal(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cacheKey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	proposalID, err, status := cdigest.ParseRequest(w, r, "proposal_id")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleDAOProposalInGroup(contract, proposalID)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			cdigest.HTTP2WriteCache(w, cacheKey, hd.expireShortLived)
		}
	}
}

func (hd *Handlers) handleDAOProposalInGroup(contract, proposalID string) (interface{}, error) {
	switch proposal, err := daodigest.DAOProposal(hd.database, contract, proposalID); {
	case err != nil:
		return nil, util.ErrNotFound.WithMessage(err, "proposal, contract %s, proposalID %s", contract, proposalID)
	case proposal == nil:
		return nil, util.ErrNotFound.Errorf("proposal, contract %s, proposalID %s", contract, proposalID)
	default:
		hal, err := hd.buildDAOProposalHal(contract, proposalID, *proposal)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildDAOProposalHal(contract, proposalID string, proposal state.ProposalStateValue) (cdigest.Hal, error) {
	h, err := hd.combineURL(daodigest.HandlerPathDAOProposal, "contract", contract, "proposal_id", proposalID)
	if err != nil {
		return nil, err
	}

	hal := cdigest.NewBaseHal(proposal, cdigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleDAODelegator(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cacheKey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	proposalID, err, status := cdigest.ParseRequest(w, r, "proposal_id")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	delegator, err, status := cdigest.ParseRequest(w, r, "address")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleDAODelegatorInGroup(contract, proposalID, delegator)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			cdigest.HTTP2WriteCache(w, cacheKey, hd.expireShortLived)
		}
	}
}

func (hd *Handlers) handleDAODelegatorInGroup(contract, proposalID, delegator string) (interface{}, error) {
	switch delegatorInfo, err := daodigest.DAODelegatorInfo(hd.database, contract, proposalID, delegator); {
	case err != nil:
		return nil, util.ErrNotFound.WithMessage(err, "delegator info, contract %s, proposalID %s, delegator %s", contract, proposalID, delegator)
	case delegatorInfo == nil:
		return nil, util.ErrNotFound.Errorf("delegator info, contract %s, proposalID %s, delegator %s", contract, proposalID, delegator)
	default:
		hal, err := hd.buildDAODelegatorHal(contract, proposalID, delegator, *delegatorInfo)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildDAODelegatorHal(
	contract, proposalID, delegator string,
	delegatorInfo types.DelegatorInfo,
) (cdigest.Hal, error) {
	h, err := hd.combineURL(
		daodigest.HandlerPathDAODelegator,
		"contract", contract,
		"proposal_id", proposalID,
		"address", delegator,
	)
	if err != nil {
		return nil, err
	}

	hal := cdigest.NewBaseHal(delegatorInfo, cdigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleDAOVoters(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cacheKey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	proposalID, err, status := cdigest.ParseRequest(w, r, "proposal_id")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleDAOVotersInGroup(contract, proposalID)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			cdigest.HTTP2WriteCache(w, cacheKey, hd.expireShortLived)
		}
	}
}

func (hd *Handlers) handleDAOVotersInGroup(contract, proposalID string) (interface{}, error) {
	switch voters, err := daodigest.DAOVoters(hd.database, contract, proposalID); {
	case err != nil:
		return nil, util.ErrNotFound.WithMessage(err, "voters, contract %s, proposalID %s", contract, proposalID)
	case voters == nil:
		return nil, util.ErrNotFound.Errorf("voters, contract %s, proposalID %s", contract, proposalID)
	default:
		hal, err := hd.buildDAOVotersHal(contract, proposalID, voters)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildDAOVotersHal(
	contract, proposalID string, voters []types.VoterInfo,
) (cdigest.Hal, error) {
	h, err := hd.combineURL(daodigest.HandlerPathDAOVoters, "contract", contract, "proposal_id", proposalID)
	if err != nil {
		return nil, err
	}

	hal := cdigest.NewBaseHal(voters, cdigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleDAOVotingPowerBox(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cacheKey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	proposalID, err, status := cdigest.ParseRequest(w, r, "proposal_id")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleDAOVotingPowerBoxInGroup(contract, proposalID)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			cdigest.HTTP2WriteCache(w, cacheKey, hd.expireShortLived)
		}
	}
}

func (hd *Handlers) handleDAOVotingPowerBoxInGroup(contract, proposalID string) (interface{}, error) {
	switch votingPowerBox, err := daodigest.DAOVotingPowerBox(hd.database, contract, proposalID); {
	case err != nil:
		return nil, util.ErrNotFound.WithMessage(err, "voting power box, contract %s, proposalID %s", contract, proposalID)
	case votingPowerBox == nil:
		return nil, util.ErrNotFound.Errorf("voting power box, contract %s, proposalID %s", contract, proposalID)

	default:
		hal, err := hd.buildDAOVotingPowerBoxHal(contract, proposalID, *votingPowerBox)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildDAOVotingPowerBoxHal(
	contract, proposalID string,
	votingPowerBox types.VotingPowerBox,
) (cdigest.Hal, error) {
	h, err := hd.combineURL(
		daodigest.HandlerPathDAOVotingPowerBox,
		"contract", contract,
		"proposal_id", proposalID,
	)
	if err != nil {
		return nil, err
	}

	hal := cdigest.NewBaseHal(votingPowerBox, cdigest.NewHalLink(h, nil))

	return hal, nil
}
