package area

import (
	"github.com/gingerxman/eel"
)

type Area struct {
	eel.RestResource
}

func (this *Area) Resource() string {
	return "area.area"
}

func (this *Area) SkipAuthCheck() bool {
	return true
}

func (r *Area) IsForDevTest() bool {
	return true
}

func (this *Area) GetParameters() map[string][]string {
	return map[string][]string{
		"GET":  []string{"name"},
	}
}

func (this *Area) Get(ctx *eel.Context) {
	req := ctx.Request
	areaName := req.GetString("name")
	
	area := eel.NewAreaService().GetAreaByName(areaName)
	
	province := map[string]interface{}{
		"id": area.Province.Id,
		"name": area.Province.Name,
	}
	city := map[string]interface{}{
		"id": area.City.Id,
		"name": area.City.Name,
	}
	district := map[string]interface{}{
		"id": area.District.Id,
		"name": area.District.Name,
	}
	
	data := map[string]interface{}{
		"province": province,
		"city": city,
		"district": district,
	}
	
	
	ctx.Response.JSON(data)
}

