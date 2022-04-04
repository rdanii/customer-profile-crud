package structs

type Users struct {
	ID           int            `json:"id" gorm:"primary_key"`
	Name         string         `json:"name"`
	Age          int            `json:"age"`
	Password     string         `json:"password"`
	Risk_profile []Risk_profile `json:"risk_profile,omitempty" gorm:"foreignkey:Userid"`
}

type Risk_profile struct {
	Userid        int     `json:"userid"`
	Mm_percent    float32 `json:"mm_percent"`
	Bond_percent  float32 `json:"bond_percent"`
	Stock_percent float32 `json:"stock_percent"`
	Total_percent float32 `json:"total_percent"`
}

type Result struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}
