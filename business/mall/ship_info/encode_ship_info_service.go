package ship_info

import (
	"context"

	"github.com/gingerxman/eel"
)

type EncodeShipInfoService struct {
	eel.ServiceBase
}

func NewEncodeShipInfoService(ctx context.Context) *EncodeShipInfoService {
	service := new(EncodeShipInfoService)
	service.Ctx = ctx
	return service
}

//Encode 对单个实体对象进行编码
func (this *EncodeShipInfoService) Encode(shipInfo *ShipInfo) *RShipInfo {
	if shipInfo == nil {
		return nil
	}

	return &RShipInfo{
		Id: shipInfo.Id,
		UserId: shipInfo.UserId,
		Name: shipInfo.Name,
		Phone: shipInfo.Phone,
		Area: shipInfo.Area,
		AreaCode: shipInfo.AreaCode,
		AreaJson: shipInfo.AreaJson,
		Address: shipInfo.Address,
		IsDefault: shipInfo.IsDefault,
		CreatedAt: shipInfo.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

//EncodeMany 对实体对象进行批量编码
func (this *EncodeShipInfoService) EncodeMany(shipInfos []*ShipInfo) []*RShipInfo {
	rDatas := make([]*RShipInfo, 0)
	for _, shipInfo := range shipInfos {
		rDatas = append(rDatas, this.Encode(shipInfo))
	}
	
	return rDatas
}

func init() {
}
