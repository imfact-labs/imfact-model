package digest

import (
	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-prescription/types"
	"github.com/ProtoconNet/mitum2/base"
	"net/http"
	"time"
)

func (hd *Handlers) handlePrescriptionDesign(w http.ResponseWriter, r *http.Request) {
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
		return hd.handlePrescriptionDesignInGroup(contract)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)

		if !shared {
			cdigest.HTTP2WriteCache(w, cacheKey, time.Second*1)
		}
	}
}

func (hd *Handlers) handlePrescriptionDesignInGroup(contract string) ([]byte, error) {
	var de types.Design
	var st base.State

	de, st, err := PrescriptionDesign(hd.database, contract)
	if err != nil {
		return nil, err
	}

	i, err := hd.buildPrescriptionDesign(contract, de, st)
	if err != nil {
		return nil, err
	}
	return hd.encoder.Marshal(i)
}

func (hd *Handlers) buildPrescriptionDesign(contract string, de types.Design, st base.State) (cdigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathPrescriptionDesign, "contract", contract)
	if err != nil {
		return nil, err
	}

	var hal cdigest.Hal
	hal = cdigest.NewBaseHal(de, cdigest.NewHalLink(h, nil))

	h, err = hd.combineURL(cdigest.HandlerPathBlockByHeight, "height", st.Height().String())
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("block", cdigest.NewHalLink(h, nil))

	for i := range st.Operations() {
		h, err := hd.combineURL(cdigest.HandlerPathOperation, "hash", st.Operations()[i].String())
		if err != nil {
			return nil, err
		}
		hal = hal.AddLink("operations", cdigest.NewHalLink(h, nil))
	}

	return hal, nil
}

func (hd *Handlers) handlePrescriptionInfo(w http.ResponseWriter, r *http.Request) {
	cacheKey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	key, err, status := cdigest.ParseRequest(w, r, "prescription_hash")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handlePrescriptionInfoInGroup(contract, key)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)

		if !shared {
			cdigest.HTTP2WriteCache(w, cacheKey, time.Second*1)
		}
	}
}

func (hd *Handlers) handlePrescriptionInfoInGroup(contract, hash string) ([]byte, error) {
	info, err := PrescriptionInfo(hd.database, contract, hash)
	if err != nil {
		return nil, err
	}

	i, err := hd.buildPrescriptionInfoHal(contract, *info)
	if err != nil {
		return nil, err
	}
	return hd.encoder.Marshal(i)
}

func (hd *Handlers) buildPrescriptionInfoHal(
	contract string, info types.PrescriptionInfo) (cdigest.Hal, error) {
	h, err := hd.combineURL(
		HandlerPathPrescriptionInfo,
		"contract", contract, "prescription_hash", info.PrescriptionHash())
	if err != nil {
		return nil, err
	}

	var hal cdigest.Hal
	hal = cdigest.NewBaseHal(
		info,
		cdigest.NewHalLink(h, nil),
	)

	return hal, nil
}
