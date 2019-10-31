package area

import (
	"github.com/gingerxman/eel"
)

type YouzanAreaList struct {
	eel.RestResource
}

func (this *YouzanAreaList) Resource() string {
	return "area.youzan_area_list"
}

func (this *YouzanAreaList) SkipAuthCheck() bool {
	return true
}

func (this *YouzanAreaList) GetParameters() map[string][]string {
	return map[string][]string{
		"GET":  []string{},
	}
}

func (this *YouzanAreaList) Get(ctx *eel.Context) {
	areaService := eel.NewAreaService()
	
	ctx.Response.JSON(areaService.GetRawData())
}

