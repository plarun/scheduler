package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/go-sql-driver/mysql"
	"github.com/plarun/scheduler/config"
)

// singleton instance for database
var dbConn *sql.DB = nil
var once sync.Once

// GetDatabase returns the singleton DB instance
func GetDatabase() *sql.DB {
	if dbConn == nil {
		once.Do(
			func() {
				ConnectDB()
			})
	}
	return dbConn
}

// ConnectDB connects to the Database
func ConnectDB() error {
	dc := config.GetDBConfig()
	addr := fmt.Sprintf("%s:%s", dc.Host, dc.Port)

	cfg := mysql.Config{
		User:      dc.User,
		Passwd:    dc.Password,
		Net:       "tcp",
		Addr:      addr,
		DBName:    dc.Schema,
		ParseTime: true,
	}

	// get db handle
	var err error
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return fmt.Errorf("unable to open connection to mysql db: %v", err)
	}
	dbConn = db

	// ping Database to check connectivity
	pingErr := db.Ping()
	if pingErr != nil {
		return fmt.Errorf("unable to ping to mysql db on %s: %v", addr, err)
	}

	log.Println("Connected to DB")
	return nil
}
