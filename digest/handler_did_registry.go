package digest

import (
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-did-registry/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"time"
)

func (hd *Handlers) handleDIDDesign(w http.ResponseWriter, r *http.Request) {
	cacheKey := currencydigest.CacheKeyPath(r)
	if err := currencydigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := currencydigest.ParseRequest(w, r, "contract")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleDIDDesignInGroup(contract)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)

		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleDIDDesignInGroup(contract string) ([]byte, error) {
	var de types.Design
	var st base.State

	de, st, err := DIDDesign(hd.database, contract)
	if err != nil {
		return nil, err
	}

	i, err := hd.buildDIDDesign(contract, de, st)
	if err != nil {
		return nil, err
	}
	return hd.encoder.Marshal(i)
}

func (hd *Handlers) buildDIDDesign(contract string, de types.Design, st base.State) (currencydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathDIDDesign, "contract", contract)
	if err != nil {
		return nil, err
	}

	var hal currencydigest.Hal
	hal = currencydigest.NewBaseHal(de, currencydigest.NewHalLink(h, nil))

	h, err = hd.combineURL(currencydigest.HandlerPathBlockByHeight, "height", st.Height().String())
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("block", currencydigest.NewHalLink(h, nil))

	for i := range st.Operations() {
		h, err := hd.combineURL(currencydigest.HandlerPathOperation, "hash", st.Operations()[i].String())
		if err != nil {
			return nil, err
		}
		hal = hal.AddLink("operations", currencydigest.NewHalLink(h, nil))
	}

	return hal, nil
}

func (hd *Handlers) handleDIDData(w http.ResponseWriter, r *http.Request) {
	cacheKey := currencydigest.CacheKeyPath(r)
	if err := currencydigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := currencydigest.ParseRequest(w, r, "contract")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	pkey, err, status := currencydigest.ParseRequest(w, r, "pubKey")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}
	pubKey := strings.TrimPrefix(pkey, "0x")
	// reform pubkey
	key := "04" + pubKey[len(pubKey)-128:]

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleDIDDataInGroup(contract, key)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)

		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleDIDDataInGroup(contract, key string) ([]byte, error) {
	data, st, err := DIDData(hd.database, contract, key)
	if err != nil {
		return nil, err
	}

	i, err := hd.buildDIDDataHal(contract, *data, st)
	if err != nil {
		return nil, err
	}
	return hd.encoder.Marshal(i)
}

func (hd *Handlers) buildDIDDataHal(
	contract string, data types.Data, st base.State) (currencydigest.Hal, error) {
	h, err := hd.combineURL(
		HandlerPathDIDData,
		"contract", contract, "pubKey", data.PubKey())
	if err != nil {
		return nil, err
	}

	var hal currencydigest.Hal
	hal = currencydigest.NewBaseHal(data, currencydigest.NewHalLink(h, nil))
	h, err = hd.combineURL(currencydigest.HandlerPathBlockByHeight, "height", st.Height().String())
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("block", currencydigest.NewHalLink(h, nil))

	for i := range st.Operations() {
		h, err := hd.combineURL(currencydigest.HandlerPathOperation, "hash", st.Operations()[i].String())
		if err != nil {
			return nil, err
		}
		hal = hal.AddLink("operations", currencydigest.NewHalLink(h, nil))
	}

	return hal, nil
}

func (hd *Handlers) handleDIDDocument(w http.ResponseWriter, r *http.Request) {
	cacheKey := currencydigest.CacheKeyPath(r)
	if err := currencydigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := currencydigest.ParseRequest(w, r, "contract")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	did := currencydigest.ParseStringQuery(r.URL.Query().Get("did"))
	if len(did) < 1 {
		currencydigest.HTTP2ProblemWithError(w, errors.Errorf("invalid DID"), http.StatusBadRequest)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleDIDDocumentInGroup(contract, did)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)

		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleDIDDocumentInGroup(contract, key string) ([]byte, error) {
	doc, st, err := DIDDocument(hd.database, contract, key)
	if err != nil {
		return nil, err
	}

	i, err := hd.buildDIDDocumentHal(contract, *doc, st)
	if err != nil {
		return nil, err
	}
	return hd.encoder.Marshal(i)
}

func (hd *Handlers) buildDIDDocumentHal(
	contract string, doc types.Document, st base.State) (currencydigest.Hal, error) {
	//h, err := hd.combineURL(
	//	HandlerPathDIDDocument,
	//	"contract", contract)
	//if err != nil {
	//	return nil, err
	//}

	var hal currencydigest.Hal
	hal = currencydigest.NewBaseHal(doc, currencydigest.NewHalLink("", nil))
	h, err := hd.combineURL(currencydigest.HandlerPathBlockByHeight, "height", st.Height().String())
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("block", currencydigest.NewHalLink(h, nil))

	for i := range st.Operations() {
		h, err := hd.combineURL(currencydigest.HandlerPathOperation, "hash", st.Operations()[i].String())
		if err != nil {
			return nil, err
		}
		hal = hal.AddLink("operations", currencydigest.NewHalLink(h, nil))
	}

	return hal, nil
}
