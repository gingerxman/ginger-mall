package material

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/eel/config"
	"time"
)

var OSS_KEY_ID string = config.ServiceConfig.String("aliyun::OSS_KEY_ID")
var OSS_KEY_SECRET string = config.ServiceConfig.String("aliyun::OSS_KEY_SECRET")
var OSS_BUCKET string = config.ServiceConfig.String("aliyun::OSS_BUCKET")
var OSS_ENDPOINT string = config.ServiceConfig.String("aliyun::OSS_ENDPOINT")
var CDN_HOST string = config.ServiceConfig.String("aliyun::CDN_HOST")

type Image struct {
	eel.RestResource
}

func (this *Image) Resource() string {
	return "material.image"
}

func (this *Image) SkipAuthCheck() bool {
	return true
}


func (this *Image) GetParameters() map[string][]string {
	return map[string][]string{
		"POST": []string{},
	}
}

func (this *Image) Post(ctx *eel.Context) {
	req := ctx.Request
	//bCtx := ctx.GetBusinessContext()
	
	filename := time.Now().Format("2006_01_02_15_04_05.00000")
	
	//TODO: 上传图片时，需要区分corp
	yunPath := fmt.Sprintf("upload/ginger_image/%d/%s.jpg", 1, filename)
	
	file, _, err := req.GetFile("image")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	
	//imageBytes, err := ioutil.ReadFile("./95.jpg")
	
	client, err := oss.New(OSS_ENDPOINT, OSS_KEY_ID, OSS_KEY_SECRET)
	if err != nil {
		eel.Logger.Error(err)
		ctx.Response.Error("resource:oss.connect_oss_fail", err.Error())
		return
	}
	
	bucket, err := client.Bucket(OSS_BUCKET)
	if err != nil {
		eel.Logger.Error(err)
		ctx.Response.Error("resource:oss.invalid_bucket", err.Error())
		return
	}
	
	err = bucket.PutObject(yunPath, file)
	if err != nil {
		eel.Logger.Error(err)
		ctx.Response.Error("resource:oss.upload_fail", err.Error())
		return
	}
	
	cdnPath := fmt.Sprintf("http://%s/%s", CDN_HOST, yunPath)
	
	ctx.Response.JSON(eel.Map{
		"path": cdnPath,
	})
}
