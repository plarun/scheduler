package allocator

import (
	"fmt"
	"log"
	"net"

	db "github.com/plarun/scheduler/internal/allocator/db/mysql"
	"github.com/plarun/scheduler/internal/allocator/handler"
	"github.com/plarun/scheduler/internal/allocator/service"
	"github.com/plarun/scheduler/proto"
	"google.golang.org/grpc"
)

func Serve(port int) {
	addr := fmt.Sprintf(":%d", port)

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to create tcp listener on %s", addr)
	}

	log.Printf("tcp listening on %s", addr)

	// connect to mysql db
	if err := db.ConnectDB(); err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Println("Allocator starting...")

		alloc := service.NewAllocator()

		if err := alloc.Start(); err != nil {
			panic(err)
		}
	}()

	ser := grpc.NewServer()

	// register all grpc services here
	proto.RegisterWaitTaskServiceServer(ser, handler.NewTaskExecService())

	if err := ser.Serve(listen); err != nil {
		log.Fatal("failed to listen")
	}
}
