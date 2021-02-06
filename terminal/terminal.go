package terminal

import (
	"io"
	"log"
	"os"
)

var buildOut io.Writer
var buildErr io.Writer
var processOut io.Writer
var processErr io.Writer

func Start() error {
	log.Println("skewer start.")
	return nil
}

func End() {
	log.Println("skewer terminated.")
	os.Exit(0)
}

func Verbose(args ...interface{}) {
	log.Println(args...)
}
