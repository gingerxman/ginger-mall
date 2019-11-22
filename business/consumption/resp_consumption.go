package consumption

import "github.com/gingerxman/ginger-mall/business/account"

type RConsumptionRecord struct {
	Id int `json:"id"`
	Money int `json:"consume_money"`
	ConsumeCount int `json:"consume_count"`
	User *account.RUser `json:"user"`
	CreatedAt string `json:"created_at"`
	LatestConsumeTime string `json:"latest_consume_time"`
}

func init() {
}
