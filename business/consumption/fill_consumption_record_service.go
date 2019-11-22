package consumption

import (
	"context"
	"github.com/gingerxman/ginger-mall/business/account"
	
	"github.com/gingerxman/eel"
)

type FillConsumptionRecordService struct {
	eel.ServiceBase
}

func NewFillConsumptionRecordService(ctx context.Context) *FillConsumptionRecordService {
	service := new(FillConsumptionRecordService)
	service.Ctx = ctx
	return service
}

func (this *FillConsumptionRecordService) Fill(records []*ConsumptionRecord, option map[string]interface{}) {
	if len(records) == 0 {
		return
	}
	
	ids := make([]int, 0)
	for _, record := range records {
		ids = append(ids, record.Id)
	}
	
	this.fillUsers(records, ids)
}

func (this *FillConsumptionRecordService) fillUsers(records []*ConsumptionRecord, ids []int) {
	if len(ids) == 0 {
		return
	}
	
	userIds := make([]int, 0)
	
	for _, record := range records {
		userIds = append(userIds, record.UserId)
	}
	
	users := account.NewUserRepository(this.Ctx).GetUsers(userIds)
	
	//构建<id, user>
	id2user := make(map[int]*account.User)
	for _, user := range users {
		id2user[user.Id] = user
	}
	
	//填充record.User
	for _, record := range records {
		if user, ok := id2user[record.UserId]; ok {
			record.User =  user
		}
	}
}

func init() {
}
