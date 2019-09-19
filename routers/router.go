package routers

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/eel/handler/rest/console"
	"github.com/gingerxman/eel/handler/rest/op"
	"github.com/gingerxman/ginger-mall/rest/dev"
	"github.com/gingerxman/ginger-mall/rest/mall"
	"github.com/gingerxman/ginger-mall/rest/product"
)

func init() {
	eel.RegisterResource(&console.Console{})
	eel.RegisterResource(&op.Health{})
	
	/*
	 product
	 */
	//category
	eel.RegisterResource(&product.Category{})
	eel.RegisterResource(&product.DisabledCategory{})
	eel.RegisterResource(&product.SubCategories{})
	//label
	eel.RegisterResource(&product.ProductLabel{})
	eel.RegisterResource(&product.ProductLabels{})
	eel.RegisterResource(&product.CorpProductLabels{})
	eel.RegisterResource(&product.DisabledCategory{})
	//property
	eel.RegisterResource(&product.ProductProperty{})
	eel.RegisterResource(&product.ProductPropertyValue{})
	eel.RegisterResource(&product.ProductProperties{})
	eel.RegisterResource(&product.CorpProductProperties{})
	//product
	eel.RegisterResource(&product.Product{})
	eel.RegisterResource(&product.OffshelfProducts{})
	eel.RegisterResource(&product.OnshelfProducts{})
	
	/*
	 mall
	 */
	eel.RegisterResource(&mall.SubCategories{})

	/*
	 dev
	 */
	eel.RegisterResource(&dev.BDDReset{})
}