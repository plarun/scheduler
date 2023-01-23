package validator

import (
	"fmt"
	"log"
	"net"

	db "github.com/plarun/scheduler/internal/validator/db/mysql"
	"github.com/plarun/scheduler/internal/validator/handler"

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
	db.ConnectDB()

	ser := grpc.NewServer()

	// register all grpc services here
	proto.RegisterValidatedActionServiceServer(ser, handler.NewJobValidationService())

	if err := ser.Serve(listen); err != nil {
		log.Fatalf("failed to listen")
	}
}
