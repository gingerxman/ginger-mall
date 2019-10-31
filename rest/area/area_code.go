package area

import (
	"github.com/gingerxman/eel"
)

type AreaCode struct {
	eel.RestResource
}

func (this *AreaCode) Resource() string {
	return "area.area_code"
}

func (this *AreaCode) SkipAuthCheck() bool {
	return true
}

func (this *AreaCode) GetParameters() map[string][]string {
	return map[string][]string{
		"GET":  []string{"name"},
	}
}

func (this *AreaCode) Get(ctx *eel.Context) {
	req := ctx.Request
	areaName := req.GetString("name")
	
	area := eel.NewAreaService().GetAreaByName(areaName)
	
	ctx.Response.JSON(area.District.Id)
}

