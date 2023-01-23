package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/plarun/scheduler/config"
	cli "github.com/plarun/scheduler/internal/client"
	tm "github.com/plarun/scheduler/pkg/time"
	"google.golang.org/grpc/status"
)

func main() {
	// export configs
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	// set logger
	setLogger()

	if len(os.Args) < 2 {
		fmt.Println("missing sub command")
		return
	}

	// new sub command
	cmd := cli.New(os.Args[1])
	if cmd == nil {
		fmt.Println("invalid sub command")
		return
	}
	log.Println("new command created")

	// parse the sub command and its flags
	if err := cmd.Parse(os.Args[2:]); err != nil {
		fmt.Printf("Error: %v", err)
		fmt.Println(cmd.Usage())
		return
	}
	log.Println("sub command parsed", cmd)

	// execute the sub command
	if err := cmd.Exec(); err != nil {
		if st, ok := status.FromError(err); ok {
			fmt.Println(st.Message())
		} else {
			fmt.Printf("Error: %v", err)
		}
	}
	log.Println("sub command executed")
}

func setLogger() {
	cfg := config.GetAppConfig()
	lg := cfg.Service.Client.Logs

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
		log.Fatal("unable to create log output file")
	}

	log.SetOutput(file)
	log.SetFlags(log.Lmicroseconds | log.Ldate | log.Lshortfile)
}
