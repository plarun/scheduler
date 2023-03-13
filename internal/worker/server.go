package worker

import (
	"log"

	"github.com/plarun/scheduler/internal/worker/service"
)

func Serve(port int) {
	log.Println("Worker starting...")

	wk := service.NewWorker(10)
	go wk.Start()

	fd := service.NewTaskFeed(wk)
	if err := fd.Start(); err != nil {
		log.Fatal(err)
	}
}
