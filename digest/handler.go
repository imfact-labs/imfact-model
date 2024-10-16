package digest

import (
	"context"
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	isaacnetwork "github.com/ProtoconNet/mitum2/isaac/network"
	"github.com/ProtoconNet/mitum2/network/quicmemberlist"
	"github.com/ProtoconNet/mitum2/network/quicstream"
	"net/http"
	"time"

	"github.com/ProtoconNet/mitum-currency/v3/digest/network"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/logging"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/singleflight"
)

var (
	HandlerPathNFTAllApproved        = `/nft/{contract:(?i)` + types.REStringAddressString + `}/account/{address:(?i)` + types.REStringAddressString + `}/allapproved` // revive:disable-line:line-length-limit
	HandlerPathNFTCollection         = `/nft/{contract:(?i)` + types.REStringAddressString + `}`
	HandlerPathNFT                   = `/nft/{contract:(?i)` + types.REStringAddressString + `}/nftidx/{nft_idx:[0-9]+}`
	HandlerPathNFTs                  = `/nft/{contract:(?i)` + types.REStringAddressString + `}/nfts`
	HandlerPathNFTCount              = `/nft/{contract:(?i)` + types.REStringAddressString + `}/totalsupply`
	HandlerPathDIDService            = `/did/{contract:(?i)` + types.REStringAddressString + `}`
	HandlerPathDIDCredential         = `/did/{contract:(?i)` + types.REStringAddressString + `}/template/{template_id:` + types.ReSpecialCh + `}/credential/{credential_id:` + types.ReSpecialCh + `}`
	HandlerPathDIDTemplate           = `/did/{contract:(?i)` + types.REStringAddressString + `}/template/{template_id:` + types.ReSpecialCh + `}`
	HandlerPathDIDCredentials        = `/did/{contract:(?i)` + types.REStringAddressString + `}/template/{template_id:` + types.ReSpecialCh + `}/credentials`
	HandlerPathDIDHolder             = `/did/{contract:(?i)` + types.REStringAddressString + `}/holder/{holder:(?i)` + types.REStringAddressString + `}` // revive:disable-line:line-length-limit
	HandlerPathTimeStampDesign       = `/timestamp/{contract:(?i)` + types.REStringAddressString + `}`
	HandlerPathTimeStampItem         = `/timestamp/{contract:(?i)` + types.REStringAddressString + `}/project/{project_id:` + types.ReSpecialCh + `}/idx/{timestamp_idx:[0-9]+}`
	HandlerPathToken                 = `/token/{contract:(?i)` + base.REStringAddressString + `}`
	HandlerPathTokenBalance          = `/token/{contract:(?i)` + base.REStringAddressString + `}/account/{address:(?i)` + base.REStringAddressString + `}` // revive:disable-line:line-length-limit
	HandlerPathPoint                 = `/point/{contract:(?i)` + base.REStringAddressString + `}`
	HandlerPathPointBalance          = `/point/{contract:(?i)` + base.REStringAddressString + `}/account/{address:(?i)` + base.REStringAddressString + `}` // revive:disable-line:line-length-limit
	HandlerPathDAOService            = `/dao/{contract:(?i)` + types.REStringAddressString + `}`
	HandlerPathDAOProposal           = `/dao/{contract:(?i)` + types.REStringAddressString + `}/proposal/{proposal_id:` + types.ReSpecialCh + `}`
	HandlerPathDAODelegator          = `/dao/{contract:(?i)` + types.REStringAddressString + `}/proposal/{proposal_id:` + types.ReSpecialCh + `}/registrant/{address:(?i)` + types.REStringAddressString + `}`
	HandlerPathDAOVoters             = `/dao/{contract:(?i)` + types.REStringAddressString + `}/proposal/{proposal_id:` + types.ReSpecialCh + `}/voter`
	HandlerPathDAOVotingPowerBox     = `/dao/{contract:(?i)` + types.REStringAddressString + `}/proposal/{proposal_id:` + types.ReSpecialCh + `}/votingpower` // revive:disable-line:line-length-limit
	HandlerPathStorageDesign         = `/storage/{contract:(?i)` + types.REStringAddressString + `}`
	HandlerPathStorageData           = `/storage/{contract:(?i)` + types.REStringAddressString + `}/datakey/{data_key:` + types.ReSpecialCh + `}`
	HandlerPathStorageHistory        = `/storage/{contract:(?i)` + types.REStringAddressString + `}/datakey/{data_key:` + types.ReSpecialCh + `}/history`
	HandlerPathPrescriptionDesign    = `/prescription/{contract:(?i)` + types.REStringAddressString + `}`
	HandlerPathPrescriptionInfo      = `/prescription/{contract:(?i)` + types.REStringAddressString + `}/hash/{prescription_hash:` + types.ReSpecialCh + `}`
	HandlerPathDIDDesign             = `/did-registry/{contract:(?i)` + types.REStringAddressString + `}`
	HandlerPathDIDData               = `/did-registry/{contract:(?i)` + types.REStringAddressString + `}/did/{pubKey:` + types.ReSpecialCh + `}`
	HandlerPathDIDDocument           = `/did-registry/{contract:(?i)` + types.REStringAddressString + `}/document`
	HandlerPathDmileDesign           = `/dmile/{contract:(?i)` + types.REStringAddressString + `}`
	HandlerPathDmileDataByTxID       = `/dmile/{contract:(?i)` + types.REStringAddressString + `}/txhash/{tx_hash:` + types.ReSpecialCh + `}`
	HandlerPathDmileDataByMerkleRoot = `/dmile/{contract:(?i)` + types.REStringAddressString + `}/merkleroot/{merkle_root:` + types.ReSpecialCh + `}`
	HandlerPathResource              = `/resource`
)

func init() {
	if b, err := currencydigest.JSON.Marshal(currencydigest.UnknownProblem); err != nil {
		panic(err)
	} else {
		currencydigest.UnknownProblemJSON = b
	}
}

type Handlers struct {
	*zerolog.Logger
	networkID       base.NetworkID
	encoders        *encoder.Encoders
	encoder         encoder.Encoder
	database        *currencydigest.Database
	cache           currencydigest.Cache
	nodeInfoHandler currencydigest.NodeInfoHandler
	send            func(interface{}) (base.Operation, error)
	client          func() (*isaacnetwork.BaseClient, *quicmemberlist.Memberlist, []quicstream.ConnInfo, error)
	router          *mux.Router
	routes          map[ /* path */ string]*mux.Route
	itemsLimiter    func(string /* request type */) int64
	rg              *singleflight.Group
	expireNotFilled time.Duration
}

func NewHandlers(
	ctx context.Context,
	networkID base.NetworkID,
	encs *encoder.Encoders,
	enc encoder.Encoder,
	st *currencydigest.Database,
	cache currencydigest.Cache,
	router *mux.Router,
	routes map[string]*mux.Route,
) *Handlers {
	var log *logging.Logging
	if err := util.LoadFromContextOK(ctx, launch.LoggingContextKey, &log); err != nil {
		return nil
	}

	return &Handlers{
		Logger:          log.Log(),
		networkID:       networkID,
		encoders:        encs,
		encoder:         enc,
		database:        st,
		cache:           cache,
		router:          router,
		routes:          routes,
		itemsLimiter:    currencydigest.DefaultItemsLimiter,
		rg:              &singleflight.Group{},
		expireNotFilled: time.Second * 1,
	}
}

func (hd *Handlers) Initialize() error {
	//cors := handlers.CORS(
	//	handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
	//	handlers.AllowedHeaders([]string{"content-type"}),
	//	handlers.AllowedOrigins([]string{"*"}),
	//	handlers.AllowCredentials(),
	//)
	//hd.router.Use(cors)

	hd.setHandlers()

	return nil
}

func (hd *Handlers) SetLimiter(f func(string) int64) *Handlers {
	hd.itemsLimiter = f

	return hd
}

func (hd *Handlers) Cache() currencydigest.Cache {
	return hd.cache
}

func (hd *Handlers) Router() *mux.Router {
	return hd.router
}

func (hd *Handlers) Handler() http.Handler {
	return network.HTTPLogHandler(hd.router, hd.Logger)
}

func (hd *Handlers) setHandlers() {
	get := 1000
	_ = hd.setHandler(HandlerPathNFTCollection, hd.handleNFTCollection, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathNFTs, hd.handleNFTs, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathNFTCount, hd.handleNFTCount, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathNFTAllApproved, hd.handleNFTOperators, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathNFT, hd.handleNFT, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathDIDService, hd.handleCredentialService, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathDIDCredentials, hd.handleCredentials, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathDIDCredential, hd.handleCredential, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathDIDHolder, hd.handleHolderCredential, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathDIDTemplate, hd.handleTemplate, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathTimeStampItem, hd.handleTimeStampItem, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathTimeStampDesign, hd.handleTimeStampDesign, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathTokenBalance, hd.handleTokenBalance, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathToken, hd.handleToken, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathPointBalance, hd.handlePointBalance, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathPoint, hd.handlePoint, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathDAOService, hd.handleDAOService, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathDAOProposal, hd.handleDAOProposal, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathDAODelegator, hd.handleDAODelegator, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathDAOVoters, hd.handleDAOVoters, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathDAOVotingPowerBox, hd.handleDAOVotingPowerBox, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathStorageData, hd.handleStorageData, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathStorageDesign, hd.handleStorageDesign, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathStorageHistory, hd.handleStorageDataHistory, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathPrescriptionInfo, hd.handlePrescriptionInfo, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathPrescriptionDesign, hd.handlePrescriptionDesign, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathDIDData, hd.handleDIDData, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathDIDDesign, hd.handleDIDDesign, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathDIDDocument, hd.handleDIDDocument, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathDmileDataByTxID, hd.handleDmileDataByTxID, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathDmileDesign, hd.handleDmileDesign, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathDmileDataByMerkleRoot, hd.handleDmileDataByMerkleRoot, true, get, get).
		Methods(http.MethodOptions, "GET")
	//_ = hd.setHandler(HandlerPathResource, hd.handleResource, true, get, get).
	//	Methods(http.MethodOptions, "GET")
}

func (hd *Handlers) setHandler(prefix string, h network.HTTPHandlerFunc, useCache bool, rps, burst int) *mux.Route {
	var handler http.Handler
	if !useCache {
		handler = http.HandlerFunc(h)
	} else {
		ch := currencydigest.NewCachedHTTPHandler(hd.cache, h)

		handler = ch
	}

	var name string
	if prefix == "" || prefix == "/" {
		name = "root"
	} else {
		name = prefix
	}

	var route *mux.Route
	if r := hd.router.Get(name); r != nil {
		route = r
	} else {
		route = hd.router.Name(name)
	}

	handler = currencydigest.RateLimiter(rps, burst)(handler)

	/*
		if rules, found := hd.rateLimit[prefix]; found {
			handler = process.NewRateLimitMiddleware(
				process.NewRateLimit(rules, limiter.Rate{Limit: -1}), // NOTE by default, unlimited
				hd.rateLimitStore,
			).Middleware(handler)

			hd.Log().Debug().Str("prefix", prefix).Msg("ratelimit middleware attached")
		}
	*/

	route = route.
		Path(prefix).
		Handler(handler)

	hd.routes[prefix] = route

	return route
}

func (hd *Handlers) combineURL(path string, pairs ...string) (string, error) {
	if n := len(pairs); n%2 != 0 {
		return "", errors.Errorf("failed to combine url; uneven pairs to combine url")
	} else if n < 1 {
		u, err := hd.routes[path].URL()
		if err != nil {
			return "", errors.Wrap(err, "failed to combine url")
		}
		return u.String(), nil
	}

	u, err := hd.routes[path].URLPath(pairs...)
	if err != nil {
		return "", errors.Wrap(err, "failed to combine url")
	}
	return u.String(), nil
}
