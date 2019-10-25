package ship_info

import (
	
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/mall/ship_info"
)

type DefaultShipInfo struct {
	eel.RestResource
}

func (this *DefaultShipInfo) Resource() string {
	return "mall.default_ship_info"
}

func (this *DefaultShipInfo) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{},
		"PUT": []string{"id:int"},
	}
}

func (this *DefaultShipInfo) Get(ctx *eel.Context) {
	bCtx := ctx.GetBusinessContext()
	user := account.GetUserFromContext(bCtx)
	repository := ship_info.NewShipInfoRepository(bCtx)
	shipInfo := repository.GetDefaultShipInfoForUser(user)
	
	if shipInfo == nil {
		ctx.Response.JSON(eel.Map{})
	} else {
		fillService := ship_info.NewFillShipInfoService(bCtx)
		fillService.Fill([]*ship_info.ShipInfo{shipInfo}, eel.FillOption{
		})
		
		encodeService := ship_info.NewEncodeShipInfoService(bCtx)
		respData := encodeService.Encode(shipInfo)
		
		ctx.Response.JSON(respData)
	}
}

func (this *DefaultShipInfo) Put(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id")

	bCtx := ctx.GetBusinessContext()
	user := account.GetUserFromContext(bCtx)
	repository := ship_info.NewShipInfoRepository(bCtx)
	shipInfo := repository.GetShipInfo(id)
	err := shipInfo.SetDefault(user)
	
	if err != nil {
		eel.Logger.Error(err)
		ctx.Response.Error("default_ship_info:create_fail", err.Error())
	} else {
	
	}

	fillService := ship_info.NewFillShipInfoService(bCtx)
	fillService.Fill([]*ship_info.ShipInfo{ shipInfo }, eel.FillOption{
	})

	encodeService := ship_info.NewEncodeShipInfoService(bCtx)
	respData := encodeService.Encode(shipInfo)
	
	ctx.Response.JSON(respData)
}

