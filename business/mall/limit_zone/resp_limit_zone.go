package limit_zone

type RLimitZone struct {
	Id int `json:"id"`
	CorpId int `json:"corp_id"`
	Name string `json:"name"`
	Provinces []*LimitProvince `json:"provinces"`
	CreatedAt string `json:"created_at"`
}


func init() {
}
