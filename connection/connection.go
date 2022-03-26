package connection

import (
	"customer-profile-crud/structs"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB  *gorm.DB
	Err error
)

func Connect() {
	DB, Err = gorm.Open("mysql", "dani:so71dCaLpTm7tK0i@/dani?charset=utf8&parseTime=True")
	// DB, Err = gorm.Open("mysql", "root:@/data_nasabah?charset=utf8&parseTime=True")

	if Err != nil {
		log.Println("Connection failed", Err)
	} else {
		log.Println("Server up and running")
	}

	DB.AutoMigrate(&structs.Users{}, &structs.Risk_profile{})
}
