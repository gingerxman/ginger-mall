package limit_zone

import (
	"context"
	"github.com/gingerxman/eel"
	"sync"
	
)

type EncodeLimitZoneService struct {
	eel.ServiceBase
	sync.RWMutex
	//id2province map[int]*RProvince
}

func NewEncodeLimitZoneService(ctx context.Context) *EncodeLimitZoneService {
	service := new(EncodeLimitZoneService)
	service.Ctx = ctx
	return service
}

//func (this *EncodeLimitZoneService) getNameProvinceMap() map[int]*RProvince {
//	if this.id2province != nil {
//		return this.id2province
//	}
//
//	this.Lock()
//	defer this.Unlock()
//
//	resp, err := eel.NewResource(this.Ctx).Get("gaia", "area.provinces", eel.Map{
//	})
//
//	if err != nil{
//		eel.Logger.Error(err)
//		panic(eel.NewBusinessError("encode_limit_zone_service:get_gaia_provinces_fail", "获取gaia的area.provinces失败"))
//	}
//
//	this.id2province = make(map[int]*RProvince)
//	data, _ := resp.Data().Map()
//	provinceDatas := data["provinces"].([]interface{})
//	for _, provinceData := range provinceDatas {
//		data := provinceData.(map[string]interface{})
//		id, err := data["id"].(json.Number).Int64()
//		if err != nil {
//			eel.Logger.Error(err)
//			continue
//		}
//
//		name := data["name"].(string)
//		this.id2province[int(id)] = &RProvince{
//			Id: int(id),
//			Name: name,
//		}
//	}
//
//	return this.id2province
//}
//
//func (this *EncodeLimitZoneService) parseDestinations(destination string) []*RProvince {
//	rProvinces := make([]*RProvince, 0)
//	if destination == "" {
//		return rProvinces
//	}
//
//	id2province := this.getNameProvinceMap()
//	items := strings.Split(destination, ",")
//	for _, item := range items {
//		id, _ := strconv.Atoi(item)
//		if province, ok := id2province[id]; ok {
//			rProvinces = append(rProvinces, province)
//		}
//	}
//	return rProvinces
//}

//Encode 对单个实体对象进行编码
func (this *EncodeLimitZoneService) Encode(limitZone *LimitZone) *RLimitZone {
	if limitZone == nil {
		return nil
	}

	return &RLimitZone{
		Id: limitZone.Id,
		CorpId: limitZone.CorpId,
		Name: limitZone.Name,
		Provinces: limitZone.Provinces,
		CreatedAt: limitZone.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

//EncodeMany 对实体对象进行批量编码
func (this *EncodeLimitZoneService) EncodeMany(limitZones []*LimitZone) []*RLimitZone {
	rDatas := make([]*RLimitZone, 0)
	for _, limitZone := range limitZones {
		rDatas = append(rDatas, this.Encode(limitZone))
	}
	
	return rDatas
}

func init() {
}
