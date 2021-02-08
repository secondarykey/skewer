package skewer

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"golang.org/x/xerrors"
)

var verbose bool

func startTerminal(v bool) error {
	log.Println("skewer start.")
	verbose = v
	return nil
}

func endTerminal() {
	log.Println("skewer terminated.")
	os.Exit(0)
}

func printVerbose(args ...interface{}) {
	if verbose {
		log.Println(args...)
	}
}

func setCommandPipe(cmd *exec.Cmd) error {

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return xerrors.Errorf("StdoutPipe error: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return xerrors.Errorf("StdoutPipe error: %w", err)
	}

	setPipe(36, stdout, os.Stdout)
	setPipe(31, stderr, os.Stderr)
	return nil
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
