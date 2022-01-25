package query

import (
	"database/sql"
	"log"
	"os"
	"sync"

	"github.com/go-sql-driver/mysql"
)

const (
	defaultTime = "1900-01-01 00:00:00"
)

// singleton instance for database
var DB *Database = nil

type Database struct {
	DB      *sql.DB
	lock    *sync.Mutex
	verbose bool
}

func GetDatabase() *Database {
	if DB == nil {
		DB = &Database{
			lock:    &sync.Mutex{},
			verbose: false,
		}
	}
	return DB
}

func (db *Database) SetVerbose() {
	db.verbose = true
}

func ConnectDB() {
	log.Println("DB connecting...")

	username := os.Getenv("DBUSER")
	password := os.Getenv("DBPASS")

	cfg := mysql.Config{
		User:   username,
		Passwd: password,
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "scheduler",
	}

	// create DB conn instance
	database := GetDatabase()

	// get db handle
	var err error
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	database.DB = db

	// ping Database to check connectivity
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Println("DB connected...")
}
