package service

import (
	"log"
	"testing"

	"github.com/plarun/scheduler/config"
	db "github.com/plarun/scheduler/internal/allocator/db/mysql"
	"github.com/plarun/scheduler/internal/allocator/db/mysql/query"
)

func init() {
	// export configs
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}
	// connect to database
	if err := db.ConnectDB(); err != nil {
		log.Fatal(err)
	}
}

func TestLockForStaging(t *testing.T) {
	if err := query.LockForStaging(); err != nil {
		t.Fatal(err)
	}
}
