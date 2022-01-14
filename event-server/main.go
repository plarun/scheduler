package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/go-sql-driver/mysql"
	pb "github.com/plarun/scheduler/event-server/data"
	"github.com/plarun/scheduler/event-server/query"
	"github.com/plarun/scheduler/event-server/service"
	"google.golang.org/grpc"
)

const port = 5555

func main() {
	// Connect to sql database
	connectDB()

	// event server service
	serve()
}

// Connect to Database
func connectDB() {
	fmt.Println("DB connecting...")

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
	database := query.GetDatabase()

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

	fmt.Println("DB Connected.")
}

func serve() {
	addr := fmt.Sprintf(":%d", port)
	// Server listens on tcp port
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	// grpc server can have arguments for unary and stream as server options
	grpcServer := grpc.NewServer()
	// register all servers here
	pb.RegisterSubmitJilServer(grpcServer, service.JilServer{Database: query.GetDatabase()})
	pb.RegisterNextJobsServer(grpcServer, service.NextJobsServer{Database: query.GetDatabase()})

	fmt.Printf("Scheduler grpc server is running at port: %d\n", port)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to start server %v", err)
	}
}
