package consumption

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/consumption"
)

type UserConsumptionRecords struct {
	eel.RestResource
}

func (this *UserConsumptionRecords) Resource() string {
	return "consumption.user_consumption_records"
}

func (this *UserConsumptionRecords) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{"?filters:json"},
	}
}


func (this *UserConsumptionRecords) Get(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()

	filters := req.GetOrmFilters()
	pageInfo := req.GetPageInfo()

	corp := account.GetCorpFromContext(bCtx)
	filters["corp_id"] = corp.Id
	records, nextPageInfo := consumption.NewConsumptionRecordRepository(bCtx).GetPagedOrders(filters, pageInfo, "-updated_at")

	consumption.NewFillConsumptionRecordService(bCtx).Fill(records, eel.Map{})
	datas := consumption.NewEncodeConsumptionRecordService(bCtx).EncodeMany(records)

	ctx.Response.JSON(eel.Map{
		"records": datas,
		"pageinfo": nextPageInfo.ToMap(),
	})
}

