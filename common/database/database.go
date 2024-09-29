package database

import (
	cc "Readee-Backend/common"
	"Readee-Backend/type/table"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	host     = "server2.bsthun.com" // or the Docker service name if running in another container
	port     = 4004                 // default PostgreSQL port
	user     = "parn"               // as defined in docker-compose.yml
	password = "parn1234"           // as defined in docker-compose.yml
	dbname   = "poc2"               // as defined in docker-compose.yml
)

func Init() {
	//dsn := fmt.Sprintf("host=%s port=%d user=%s "+
	//	"password=%s dbname=%s sslmode=disable",
	//	host, port, user, password, dbname)
	dsn2 := "parn:parn1234@tcp(server2.bsthun.com:4004)/poc2"
	db, err := gorm.Open(mysql.Open(dsn2), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	} else {
		log.Println("connect to database success")
	}

	// ลบตารางทั้งหมด
	//var tables []string
	//db.Raw("SHOW TABLES").Scan(&tables)
	//for _, table := range tables {
	//	err := db.Migrator().DropTable(table)
	//	if err != nil {
	//		return
	//	}
	//}

	//db.Migrator().DropTable(&table.Log{}) // delete table
	//db.Migrator().DropTable("user_genres") // delete joined table

	cc.DB = db
	cc.DB.AutoMigrate(&table.User{})
	//cc.DB.AutoMigrate(&table.Reserve{})
	cc.DB.AutoMigrate(&table.Book{})
	cc.DB.AutoMigrate(&table.Match{})
	cc.DB.AutoMigrate(&table.Log{})
	cc.DB.AutoMigrate(&table.History{})
	cc.DB.AutoMigrate(&table.Review{})
	cc.DB.AutoMigrate(&table.Rating{})
	cc.DB.AutoMigrate(&table.Room{})
	cc.DB.AutoMigrate(&table.Message{})
	cc.DB.AutoMigrate(&table.Genre{})
}
