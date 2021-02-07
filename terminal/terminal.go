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

var verbose bool

func Start(v bool) error {
	log.Println("skewer start.")
	verbose = v
	return nil
}

func End() {
	log.Println("skewer terminated.")
	os.Exit(0)
}

func SetVerbose(v bool) {
}

func Verbose(args ...interface{}) {
	if verbose {
		log.Println(args...)
	}
}
