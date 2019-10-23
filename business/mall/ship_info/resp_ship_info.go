package ship_info

type RShipInfo struct {
	Id int `json:"id"`
	UserId int `json:"user_id"`
	Name string `json:"name"`
	Phone string `json:"phone"`
	Area string `json:"area"`
	AreaCode string `json:"area_code"`
	AreaJson string `json:"area_json"`
	Address string `json:"address"`
	IsDefault bool `json:"is_default"`
	CreatedAt string `json:"created_at"`
}


func init() {
}
