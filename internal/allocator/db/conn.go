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
var DB *Database = nil
var once sync.Once

type Database struct {
	DB   *sql.DB
	lock *sync.RWMutex
}

// GetDatabase returns the singleton DB instance
func GetDatabase() *Database {
	if DB == nil {
		once.Do(
			func() {
				DB = &Database{
					lock: &sync.RWMutex{},
				}
			})
	}
	return DB
}

// ConnectDB connects to the Database
func ConnectDB() error {
	dc := config.GetDBConfig()
	addr := fmt.Sprintf("%s:%s", dc.Host, dc.Port)

	cfg := mysql.Config{
		User:   dc.User,
		Passwd: dc.Password,
		Net:    "tcp",
		Addr:   addr,
		DBName: dc.Schema,
	}

	// create DB conn instance
	database := GetDatabase()

	// get db handle
	var err error
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return fmt.Errorf("unable to open connection to mysql db: %v", err)
	}
	database.DB = db

	// ping Database to check connectivity
	pingErr := db.Ping()
	if pingErr != nil {
		return fmt.Errorf("unable to ping to mysql db on %s: %v", addr, err)
	}

	log.Println("Connected to DB")
	return nil
}

func (db *Database) RLock() {
	db.lock.RLock()
}

func (db *Database) RUnlock() {
	db.lock.RUnlock()
}

func (db *Database) Lock() {
	db.lock.Lock()
}

func (db *Database) Unlock() {
	db.lock.Unlock()
}
