package initializers

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var MyDB *gorm.DB

func ConnectToDb() {
	var err error
	dsn := "host=localhost user=postgres password=dkSp456dxcD dbname=jwt_test port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	MyDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
}
