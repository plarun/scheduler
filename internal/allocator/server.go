package allocator

import (
	"log"

	qry "github.com/plarun/scheduler/internal/allocator/config"
	"github.com/plarun/scheduler/internal/allocator/db/mysql"
)

func Serve(port int) {
	// addr := fmt.Sprintf(":%d", port)

	// listen, err := net.Listen("tcp", addr)
	// if err != nil {
	// 	log.Fatalf("failed to create tcp listener on %s", addr)
	// }

	// log.Printf("tcp listening on %s", addr)

	// connect to mysql db
	if err := mysql.ConnectDB(); err != nil {
		log.Fatal(err)
	}

	if err := qry.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	// ser := api.RegisterGrpcServers()

	// if err := ser.Serve(listen); err != nil {
	// 	log.Fatal("failed to listen")
	// }
}
