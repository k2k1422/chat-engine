package Database

import (
	"database/sql"
	"fmt"
	"log"
	"messaging/Model"
	"messaging/Utils"
	"os"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/driver/postgres"
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

	Utils.WaitForPort(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), 30*time.Hour)

	// db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
	err = db.AutoMigrate(&Model.User{}, &Model.Chat{})
	if err != nil {
		log.Fatalf("Error auto-migrating schema: %v", err)
	}

	log.Println("Database schema migrated successfully")
}
