package structs

type Users struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Password string `json:"password"`
}

type Risk_profile struct {
	Userid        int
	Users         Users   `gorm:"foreignkey:Userid"`
	Mm_percent    float32 `json:"mm_percent"`
	Bond_percent  float32 `json:"bond_percent"`
	Stock_percent float32 `json:"stock_percent"`
	Total_percent float32 `json:"total_percent"`
}

type Result struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}
