package consumption

import (
	"context"
	"github.com/gingerxman/ginger-mall/business/account"
	
	"github.com/gingerxman/eel"
)

type EncodeConsumptionRecordService struct {
	eel.ServiceBase
}

func NewEncodeConsumptionRecordService(ctx context.Context) *EncodeConsumptionRecordService {
	service := new(EncodeConsumptionRecordService)
	service.Ctx = ctx
	return service
}

//Encode 对单个实体对象进行编码
func (this *EncodeConsumptionRecordService) Encode(record *ConsumptionRecord) *RConsumptionRecord {
	if record == nil {
		return nil
	}
	
	var rUser *account.RUser
	if record.User != nil {
		user := record.User
		rUser = &account.RUser{
			Id: user.Id,
			Name: user.Name,
			Avatar: user.Avatar,
			Sex: user.Sex,
			Code: user.Code,
		}
	}

	return &RConsumptionRecord{
		Id: record.Id,
		Money: record.Money,
		ConsumeCount: record.ConsumeCount,
		User: rUser,
		CreatedAt: record.CreatedAt.Format("2006-01-02 15:04:05"),
		LatestConsumeTime: record.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

//EncodeMany 对实体对象进行批量编码
func (this *EncodeConsumptionRecordService) EncodeMany(records []*ConsumptionRecord) []*RConsumptionRecord {
	rDatas := make([]*RConsumptionRecord, 0)
	for _, record := range records {
		rDatas = append(rDatas, this.Encode(record))
	}
	
	return rDatas
}

func init() {
}
