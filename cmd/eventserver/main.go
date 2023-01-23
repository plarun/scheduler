package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/plarun/scheduler/config"
	es "github.com/plarun/scheduler/internal/eventserver"
	tm "github.com/plarun/scheduler/pkg/time"
)

func main() {
	// export configs
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	// set logger
	setLogger()

	es.Serve(config.GetAppConfig().Service.EventServer.Port)
}

func setLogger() {
	cfg := config.GetAppConfig()
	lg := cfg.Service.EventServer.Logs

	lay, ok := tm.GetLayout(lg.DateFormat)
	if !ok {
		lay, _ = tm.GetLayout(tm.GetDefaultDateLayout())
	}

	dtFmt := time.Now().Format(lay)
	dname := fmt.Sprintf("%s/%s", cfg.AppRoot, lg.Path)
	fname := fmt.Sprintf("%s%s.%s", lg.Prefix, dtFmt, lg.Extension)
	lfile := fmt.Sprintf("%s/%s", dname, fname)

	file, err := os.OpenFile(lfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("unable to create log output file: %v", err)
	}

	log.SetOutput(file)
	log.SetFlags(log.Lmicroseconds | log.Ldate | log.Lshortfile)
}
