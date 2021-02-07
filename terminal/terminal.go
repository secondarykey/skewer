package terminal

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

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

func Verbose(args ...interface{}) {
	if verbose {
		log.Println(args...)
	}
}

func SetPipe(stdout, stderr io.ReadCloser) {
	setPipe(36, stdout, os.Stderr)
	setPipe(31, stderr, os.Stderr)
}

func setPipe(c int, r io.ReadCloser, w io.Writer) {
	go func() {
		s := bufio.NewScanner(r)
		for s.Scan() {
			fmt.Fprintf(w, "\x1b[%dm| %s\x1b[0m\n", c, s.Text())
		}
	}()
	return
}
