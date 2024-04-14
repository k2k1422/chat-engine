package Database

import (
	"database/sql"
	"log"
	"messaging/Model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Declaring the Database client connection
var Connection *sql.DB
var DBConnection *gorm.DB

const (
	host     = "localhost"
	port     = 8801
	user     = "fammy"
	password = "password"
	dbname   = "fammydb"
)

func init() {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	Sql, err := db.DB()
	if err != nil {
		log.Fatalf("Error getting database connection: %v", err)
	}

	Connection = Sql
	DBConnection = db
	// Auto-migrate the schema
	err = db.AutoMigrate(&Model.User{})
	if err != nil {
		log.Fatalf("Error auto-migrating schema: %v", err)
	}

	log.Println("Database schema migrated successfully")
}
