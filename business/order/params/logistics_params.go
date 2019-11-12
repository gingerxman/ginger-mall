package params

import (
	"encoding/json"
	
	"github.com/gingerxman/eel"
)

type LogisticsParams struct {
	Bid string `json:"bid"`
	EnableLogistics bool `json:"enable_logistics"`
	ExpressCompanyName string `json:"express_company_name"`
	ExpressNumber string `json:"express_number"`
	Shipper string `json:"shipper"`
}

func ParseLogisticsParams(shipInfo string) *LogisticsParams{
	logisticsParams := new(LogisticsParams)
	if shipInfo != ""{
		financeInfoBytes := []byte(shipInfo)
		err := json.Unmarshal(financeInfoBytes, logisticsParams)
		if err != nil{
			eel.Logger.Error(err)
			panic(eel.NewBusinessError("logistics_params:parse_logistics_info_failed", "解析物流数据失败"))
		}
	}
	return logisticsParams
}