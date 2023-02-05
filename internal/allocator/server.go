package allocator

import (
	"log"
	"time"

	db "github.com/plarun/scheduler/internal/allocator/db/mysql"
	"github.com/plarun/scheduler/internal/allocator/service"
)

const (
	SCHEDULE_CYCLE time.Duration = time.Second * 3
)

func Serve(port int) {
	log.Println("Allocator starting...")
	// connect to mysql db
	if err := db.ConnectDB(); err != nil {
		log.Fatal(err)
	}

	alloc := service.NewAllocator()

	if err := alloc.Start(); err != nil {
		panic(err)
	}
}
