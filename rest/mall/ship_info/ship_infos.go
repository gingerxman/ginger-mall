package ship_info

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/mall/ship_info"
)

type ShipInfos struct {
	eel.RestResource
}

func (this *ShipInfos) Resource() string {
	return "mall.ship_infos"
}

func (this *ShipInfos) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{},
	}
}

func (this *ShipInfos) Get(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()
	page := req.GetPageInfo()
	page.CountPerPage = 100
	filters := req.GetOrmFilters()
	repository := ship_info.NewShipInfoRepository(bCtx)
	user := account.GetUserFromContext(bCtx)
	shipInfos, _ := repository.GetAllShipInfosForUser(user, page, filters)

	fillService := ship_info.NewFillShipInfoService(bCtx)
	fillService.Fill(shipInfos, eel.FillOption{
	})

	encodeService := ship_info.NewEncodeShipInfoService(bCtx)
	rows := encodeService.EncodeMany(shipInfos)
	
	ctx.Response.JSON(eel.Map{
		"ship_infos": rows,
	})
}
