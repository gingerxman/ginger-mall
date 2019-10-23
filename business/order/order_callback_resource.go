package order

import (
	"github.com/gingerxman/eel/rest_client"
	
	"github.com/gingerxman/eel"
	"strings"
)

type OrderCallbackResource struct{
	rawResource string
	serviceName string
	resource string
	reqData eel.Map
	orderBid string
}

// SetReqData 设置请求参数
func (this *OrderCallbackResource) SetReqData(data eel.Map){
	this.reqData = data
	this.orderBid = data["bid"].(string)
}

func (this *OrderCallbackResource) DoRequest(resource *rest_client.Resource) error{

	var err error
	errCode := ""

	resp, err := resource.Put(this.serviceName, this.resource, this.reqData)

	if err != nil{
		errCode = "网络通信失败"
		if resp != nil{
			errCode, _ = resp.RespData.Get("errCode").String()
			errMsg, _ := resp.RespData.Get("errMsg").String()
			err = eel.NewBusinessError(errCode, errMsg)
		}
		eel.Logger.Error(err)
	}

	if errCode != ""{
		// TODO: 钉钉消息
		//dingMsg := fmt.Sprintf("> 订单(%s)回调(%s)失败: \n\n 错误信息: %s \n\n", this.orderBid, this.rawResource, errCode)
		//vanilla.NewDingBot().Use("xiuer").Error(dingMsg)
	}

	return err
}

func NewOrderCallbackResource(str string) *OrderCallbackResource{
	inst := new(OrderCallbackResource)
	inst.rawResource = str
	splits := strings.Split(str, ":")
	inst.serviceName = splits[0]
	inst.resource = splits[1]
	return inst
}

