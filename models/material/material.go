package material

import (
	"github.com/gingerxman/eel"
)

type Image struct {
	eel.Model
	CorpId int `gorm:"index"`
	Url string `gorm:"size:1024"`
	Width int
	Height int
}
func (this *Image) TableName() string {
	return "material_image"
}

type Video struct {
	eel.Model
	CorpId int `gorm:"index"`
	Url string `gorm:"size:1024"`
	Time float64
}
func (this *Video) TableName() string {
	return "material_video"
}





func init() {
	eel.RegisterModel(new(Image))
	eel.RegisterModel(new(Video))
}
