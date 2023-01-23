package client

import (
	"log"
	"testing"

	"github.com/plarun/scheduler/config"
)

func TestSchdDefCommand(t *testing.T) {
	// export configs
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	// set logger
	// setLogger()

	// new sub command
	cmd := New("schd_def")
	if cmd == nil {
		t.Fatal("invalid command")
	}

	// test files
	files := []string{
		"/root/go/src/github.com/plarun/scheduler/internal/client/etc/test/def/def1.txt",
		"/root/go/src/github.com/plarun/scheduler/internal/client/etc/test/def/def2.txt",
		"/root/go/src/github.com/plarun/scheduler/internal/client/etc/test/def/def3.txt",
		"/root/go/src/github.com/plarun/scheduler/internal/client/etc/test/def/def4.txt",
		"/root/go/src/github.com/plarun/scheduler/internal/client/etc/test/def/def5.txt",
		"/root/go/src/github.com/plarun/scheduler/internal/client/etc/test/def/def6.txt",
	}

	// parse the sub command and its flags
	for _, f := range files {
		args := []string{"-c", "-f", f}
		if err := cmd.Parse(args); err != nil {
			t.Fatal(err)
			return
		}

		// execute the sub command
		if err := cmd.Exec(); err != nil {
			t.Fatal(err)
		}
	}

}
