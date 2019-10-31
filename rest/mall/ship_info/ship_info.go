package ship_info

import (
	"github.com/gingerxman/eel"
	
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/mall/ship_info"
)

type ShipInfo struct {
	eel.RestResource
}

func (this *ShipInfo) Resource() string {
	return "mall.ship_info"
}

func (this *ShipInfo) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{"id:int"},
		"PUT": []string{
			"name:string",
			"phone:string",
			"area_code:string",
			"address:string",
		},
		"POST": []string{
			"id:int",
			"name:string",
			"phone:string",
			"area_code:string",
			"address:string",
		},
		"DELETE": []string{"id:int"},
	}
}

func (this *ShipInfo) Get(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id")

	bCtx := ctx.GetBusinessContext()
	repository := ship_info.NewShipInfoRepository(bCtx)
	shipInfo := repository.GetShipInfo(id)
	
	if shipInfo == nil {
		ctx.Response.Error("no_ship_info", "no_ship_info")
	} else {
		fillService := ship_info.NewFillShipInfoService(bCtx)
		fillService.Fill([]*ship_info.ShipInfo{shipInfo}, eel.FillOption{
		})
		
		encodeService := ship_info.NewEncodeShipInfoService(bCtx)
		respData := encodeService.Encode(shipInfo)
		
		ctx.Response.JSON(respData)
	}
}

func (this *ShipInfo) Put(ctx *eel.Context) {
	req := ctx.Request
	name := req.GetString("name")
	phone := req.GetString("phone")
	areaCode := req.GetString("area_code")
	address := req.GetString("address")
	
	bCtx := ctx.GetBusinessContext()
	user := account.GetUserFromContext(bCtx)
	shipInfo := ship_info.NewShipInfo(
		bCtx,
		user,
		name,
		phone,
		areaCode,
		address,
	)

	ctx.Response.JSON(eel.Map{
		"id": shipInfo.Id,
	})
}

func (this *ShipInfo) Post(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id")
	name := req.GetString("name")
	phone := req.GetString("phone")
	areaCode := req.GetString("area_code")
	address := req.GetString("address")

	bCtx := ctx.GetBusinessContext()
	user := account.GetUserFromContext(bCtx)
	repository := ship_info.NewShipInfoRepository(bCtx)
	shipInfo := repository.GetShipInfoForUser(user, id)

	_ = shipInfo.Update(
		name,
		phone,
		areaCode,
		address,
	)
	
	ctx.Response.JSON(eel.Map{
	})
}

func (this *ShipInfo) Delete(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id")

	bCtx := ctx.GetBusinessContext()
	user := account.GetUserFromContext(bCtx)
	repository := ship_info.NewShipInfoRepository(bCtx)
	shipInfo := repository.GetShipInfoForUser(user, id)
	err := shipInfo.Delete()

	if err != nil {
		eel.Logger.Error(err)
		ctx.Response.Error("shipInfo:delete_fail", err.Error())
	} else {
		ctx.Response.JSON(eel.Map{
		})
	}
}
