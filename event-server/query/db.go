package query

import (
	"database/sql"
	"sync"
)

// singleton instance for database
var DB *Database = nil

type Database struct {
	DB   *sql.DB
	lock *sync.Mutex
}

func GetDatabase() *Database {
	if DB == nil {
		return &Database{}
	}
	return DB
}
