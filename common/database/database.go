package database

import (
	cc "Readee-Backend/common"
	"Readee-Backend/type/table"
	"log"

	//myTypes"Readee-Backend/endpoint"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	// connection name: parn
	// server address: server2.bsthun.com
	// port: 4004
	// database name: poc2
	// username: parn
	// password: parn1234
	dsn2 := "parn:parn1234@tcp(server2.bsthun.com:4004)/poc2?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn2), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	} else {
		log.Println("connect to database success")

		DB = db
	}

	// Drop all tables
	// var tables []string
	// db.Raw("SHOW TABLES").Scan(&tables)
	// for _, table := range tables {
	// 	err := db.Migrator().DropTable(table)
	// 	if err != nil {
	// 		return
	// 	}
	// }

	// Drop specific table
	// db.Migrator().DropTable(&table.Log{}) // delete table
	// db.Migrator().DropTable("user_genres") // delete joined table
	// db.Exec("ALTER TABLE logs DROP COLUMN user_like_id")

	cc.DB = db
	cc.DB.AutoMigrate(&table.User{})
	cc.DB.AutoMigrate(&table.Genre{})
	cc.DB.AutoMigrate(&table.Book{})
	cc.DB.AutoMigrate(&table.Log{})
	cc.DB.AutoMigrate(&table.Room{})
	cc.DB.AutoMigrate(&table.Message{})
	cc.DB.AutoMigrate(&table.Match{})
	cc.DB.AutoMigrate(&table.History{})
	cc.DB.AutoMigrate(&table.Rating{})
	cc.DB.AutoMigrate(&table.Review{})
}
