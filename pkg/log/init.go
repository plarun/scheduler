package log

import (
	"log"
	"os"
	"time"
)

func InitLog(dir, prefix, dtfmt, ext string) {
	fname := dir + prefix + time.Now().Format(dtfmt) + "." + ext

	file, err := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("unable to create log output file")
	}

	log.SetOutput(file)
	log.SetFlags(log.Lmicroseconds | log.Ldate | log.Lshortfile)
}
