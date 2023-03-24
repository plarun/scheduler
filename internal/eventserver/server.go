package eventserver

import (
	"fmt"
	"log"
	"net"

	"github.com/plarun/scheduler/internal/eventserver/db"
	"github.com/plarun/scheduler/internal/eventserver/handler"
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

	ser := grpc.NewServer()

	// register all grpc services here
	proto.RegisterParsedActionServiceServer(ser, handler.NewParsedActionService())
	proto.RegisterTaskExecServiceServer(ser, handler.NewTaskExecService())
	proto.RegisterTaskServiceServer(ser, handler.NewTaskService())

	if err := ser.Serve(listen); err != nil {
		log.Fatal("failed to listen")
	}
}
