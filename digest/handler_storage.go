package digest

import (
	"net/http"
	"strconv"
	"time"

	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	sdigest "github.com/ProtoconNet/mitum-storage/digest"
	"github.com/ProtoconNet/mitum-storage/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

func (hd *Handlers) handleStorageDesign(w http.ResponseWriter, r *http.Request) {
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
		return hd.handleStorageDesignInGroup(contract)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)

		if !shared {
			cdigest.HTTP2WriteCache(w, cacheKey, hd.expireShortLived)
		}
	}
}

func (hd *Handlers) handleStorageDesignInGroup(contract string) ([]byte, error) {
	var de types.Design
	var st base.State

	de, st, err := sdigest.StorageDesign(hd.database, contract)
	if err != nil {
		return nil, err
	}

	i, err := hd.buildStorageDesign(contract, de, st)
	if err != nil {
		return nil, err
	}
	return hd.encoder.Marshal(i)
}

func (hd *Handlers) buildStorageDesign(contract string, de types.Design, st base.State) (cdigest.Hal, error) {
	h, err := hd.combineURL(sdigest.HandlerPathStorageDesign, "contract", contract)
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

func (hd *Handlers) handleStorageData(w http.ResponseWriter, r *http.Request) {
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

	key, err, status := cdigest.ParseRequest(w, r, "data_key")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleStorageDataInGroup(contract, key)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)

		if !shared {
			cdigest.HTTP2WriteCache(w, cacheKey, hd.expireShortLived)
		}
	}
}

func (hd *Handlers) handleStorageDataInGroup(contract, key string) ([]byte, error) {
	data, height, operation, timestamp, deleted, err := sdigest.StorageData(hd.database, contract, key)
	if err != nil {
		return nil, err
	}

	i, err := hd.buildStorageDataHal(contract, *data, height, operation, timestamp, deleted)
	if err != nil {
		return nil, err
	}
	return hd.encoder.Marshal(i)
}

func (hd *Handlers) buildStorageDataHal(
	contract string, data types.Data, height int64, operation, timestamp string, deleted bool) (cdigest.Hal, error) {
	h, err := hd.combineURL(
		sdigest.HandlerPathStorageData,
		"contract", contract, "data_key", data.DataKey())
	if err != nil {
		return nil, err
	}

	var hal cdigest.Hal
	hal = cdigest.NewBaseHal(
		struct {
			Data      types.Data `json:"data"`
			Height    int64      `json:"height"`
			Operation string     `json:"operation"`
			Timestamp string     `json:"timestamp"`
		}{Data: data, Height: height, Operation: operation, Timestamp: timestamp},
		cdigest.NewHalLink(h, nil),
	)

	h, err = hd.combineURL(cdigest.HandlerPathBlockByHeight, "height", strconv.FormatInt(height, 10))
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("block", cdigest.NewHalLink(h, nil))

	h, err = hd.combineURL(cdigest.HandlerPathOperation, "hash", operation)
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("operation", cdigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleStorageDataHistory(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	limit := cdigest.ParseLimitQuery(r.URL.Query().Get("limit"))
	offset := cdigest.ParseStringQuery(r.URL.Query().Get("offset"))
	reverse := cdigest.ParseBoolQuery(r.URL.Query().Get("reverse"))

	cacheKey := cdigest.CacheKey(
		r.URL.Path, cdigest.StringOffsetQuery(offset),
		cdigest.StringBoolQuery("reverse", reverse),
	)

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	key, err, status := cdigest.ParseRequest(w, r, "data_key")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		i, filled, err := hd.handleStorageDataHistoryInGroup(contract, key, offset, reverse, limit)

		return []interface{}{i, filled}, err
	})

	if err != nil {
		hd.Log().Err(err).Str("Issuer", contract).Msg("failed to get credentials")
		cdigest.HTTP2HandleError(w, err)

		return
	}

	var b []byte
	var filled bool
	{
		l := v.([]interface{})
		b = l[0].([]byte)
		filled = l[1].(bool)
	}

	cdigest.HTTP2WriteHalBytes(hd.encoder, w, b, http.StatusOK)

	if !shared {
		expire := hd.expireNotFilled
		if len(offset) > 0 && filled {
			expire = time.Minute
		}

		cdigest.HTTP2WriteCache(w, cacheKey, expire)
	}
}

func (hd *Handlers) handleStorageDataHistoryInGroup(
	contract, key string,
	offset string,
	reverse bool,
	l int64,
) ([]byte, bool, error) {
	var limit int64
	if l < 0 {
		limit = hd.itemsLimiter("service-credentials")
	} else {
		limit = l
	}

	var vas []cdigest.Hal
	if err := sdigest.SotrageDataHistoryByDataKey(
		hd.database, contract, key, reverse, offset, limit,
		func(data *types.Data, height int64, operation, timestamp string, deleted bool) (bool, error) {
			hal, err := hd.buildStorageDataHal(contract, *data, height, operation, timestamp, deleted)
			if err != nil {
				return false, err
			}
			vas = append(vas, hal)

			return true, nil
		},
	); err != nil {
		return nil, false, util.ErrNotFound.WithMessage(err, "data history by contract %s, data key %s", contract, key)
	} else if len(vas) < 1 {
		return nil, false, util.ErrNotFound.Errorf("data history by contract %s, data key %s", contract, key)
	}

	i, err := hd.buildStorageDataHistoryHal(contract, key, vas, offset, reverse)
	if err != nil {
		return nil, false, err
	}

	b, err := hd.encoder.Marshal(i)
	return b, int64(len(vas)) == limit, err
}

func (hd *Handlers) buildStorageDataHistoryHal(
	contract, key string,
	vas []cdigest.Hal,
	offset string,
	reverse bool,
) (cdigest.Hal, error) {
	baseSelf, err := hd.combineURL(
		sdigest.HandlerPathStorageDataHistory,
		"contract", contract,
		"data_key", key,
	)
	if err != nil {
		return nil, err
	}

	self := baseSelf
	if len(offset) > 0 {
		self = cdigest.AddQueryValue(baseSelf, cdigest.StringOffsetQuery(offset))
	}
	if reverse {
		self = cdigest.AddQueryValue(baseSelf, cdigest.StringBoolQuery("reverse", reverse))
	}

	var hal cdigest.Hal
	hal = cdigest.NewBaseHal(vas, cdigest.NewHalLink(self, nil))

	h, err := hd.combineURL(sdigest.HandlerPathStorageDesign, "contract", contract)
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("service", cdigest.NewHalLink(h, nil))

	var nextOffset string

	if len(vas) > 0 {
		va, ok := vas[len(vas)-1].Interface().(struct {
			Data      types.Data `json:"data"`
			Height    int64      `json:"height"`
			Operation string     `json:"operation"`
			Timestamp string     `json:"timestamp"`
		})
		if !ok {
			return nil, errors.Errorf("failed to build storage data history hal")
		}
		nextOffset = strconv.FormatInt(va.Height, 10)
	}

	if len(nextOffset) > 0 {
		next := baseSelf
		next = cdigest.AddQueryValue(next, cdigest.StringOffsetQuery(nextOffset))

		if reverse {
			next = cdigest.AddQueryValue(next, cdigest.StringBoolQuery("reverse", reverse))
		}

		hal = hal.AddLink("next", cdigest.NewHalLink(next, nil))
	}

	hal = hal.AddLink("reverse", cdigest.NewHalLink(cdigest.AddQueryValue(baseSelf, cdigest.StringBoolQuery("reverse", !reverse)), nil))

	return hal, nil
}

func (hd *Handlers) handleStorageDataCount(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cachekey := cdigest.CacheKey(
		r.URL.Path,
	)

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	deleted := cdigest.ParseBoolQuery(r.URL.Query().Get("deleted"))

	v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		i, err := hd.handleStorageDataCountInGroup(contract, deleted)

		return i, err
	})

	if err != nil {
		hd.Log().Err(err).Str("contract", contract).Msg("failed to count nft")
		cdigest.HTTP2HandleError(w, err)

		return
	}

	cdigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)

	if !shared {
		expire := hd.expireNotFilled
		cdigest.HTTP2WriteCache(w, cachekey, expire)
	}
}

func (hd *Handlers) handleStorageDataCountInGroup(
	contract string, deleted bool,
) ([]byte, error) {
	count, err := sdigest.DataCountByContract(
		hd.database, contract, deleted,
	)
	if err != nil {
		return nil, util.ErrNotFound.WithMessage(err, "data count by contract, %s", contract)
	}

	i, err := hd.buildStorageDataCountHal(contract, count)
	if err != nil {
		return nil, err
	}

	b, err := hd.encoder.Marshal(i)
	return b, err
}

func (hd *Handlers) buildStorageDataCountHal(
	contract string,
	count int64,
) (cdigest.Hal, error) {
	baseSelf, err := hd.combineURL(sdigest.HandlerPathStorageDataCount, "contract", contract)
	if err != nil {
		return nil, err
	}

	self := baseSelf

	var m struct {
		Contract  string `json:"contract"`
		DataCount int64  `json:"data_count"`
	}

	m.Contract = contract
	m.DataCount = count

	var hal cdigest.Hal
	hal = cdigest.NewBaseHal(m, cdigest.NewHalLink(self, nil))

	h, err := hd.combineURL(sdigest.HandlerPathStorageDesign, "contract", contract)
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("collection", cdigest.NewHalLink(h, nil))

	return hal, nil
}
