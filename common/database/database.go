package database

import (
	cc "Readee-Backend/common"
	"Readee-Backend/type/table"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() {
	dsn := "host=server2.bsthun.com user=readee password=Readee_1234 dbname=readee port=4001 sslmode=disable TimeZone=Asia/Bangkok"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	cc.DB = db
	cc.DB.AutoMigrate(&table.User{})
	//cc.DB.AutoMigrate(&table.Reserve{})
	cc.DB.AutoMigrate(&table.Book{})
	//cc.DB.AutoMigrate(&table.Match{})
	//cc.DB.AutoMigrate(&table.Buddy{})
	//cc.DB.AutoMigrate(&table.History{})
	//cc.DB.AutoMigrate(&table.Review{})
	//cc.DB.AutoMigrate(&table.Rating{})
	//cc.DB.AutoMigrate(&table.Room{})
	//cc.DB.AutoMigrate(&table.Message{})
}
