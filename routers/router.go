package routers

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/eel/handler/rest/console"
	"github.com/gingerxman/eel/handler/rest/op"
	"github.com/gingerxman/ginger-mall/rest/dev"
	"github.com/gingerxman/ginger-mall/rest/product"
)

func init() {
	eel.RegisterResource(&console.Console{})
	eel.RegisterResource(&op.Health{})
	
	eel.RegisterResource(&product.Category{})
	eel.RegisterResource(&product.SubCategories{})

	eel.RegisterResource(&dev.BDDReset{})
}