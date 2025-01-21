package digest

import (
	"net/http"
	"time"

	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/shirou/gopsutil/mem"
)

func (hd *Handlers) handleResource(w http.ResponseWriter, r *http.Request) {
	cacheKey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleResourceInGroup()
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			cdigest.HTTP2WriteCache(w, cacheKey, time.Millisecond*500)
		}
	}
}

func (hd *Handlers) handleResourceInGroup() (interface{}, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	//partitions, err := disk.Partitions(true)
	//if err != nil {
	//	return nil, err
	//}

	//var diskUsage []disk.UsageStat
	//
	//for _, partition := range partitions {
	//	usage, err := disk.Usage(partition.Mountpoint)
	//	if err != nil {
	//		return nil, err
	//	}
	//	diskUsage = append(diskUsage, *usage)
	//}

	var m struct {
		MemInfo mem.VirtualMemoryStat `json:"mem"`
		//DiskUsage []disk.UsageStat      `json:"disk"`
	}

	m.MemInfo = *memInfo
	//m.DiskUsage = diskUsage

	hal, err := hd.buildResourceHal(m)
	if err != nil {
		return nil, err
	}
	return hd.encoder.Marshal(hal)

}

func (hd *Handlers) buildResourceHal(resource interface{}) (cdigest.Hal, error) {
	hal := cdigest.NewBaseHal(resource, cdigest.NewHalLink(HandlerPathResource, nil))

	return hal, nil
}
