package skewer

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/xerrors"
)

var verbose bool

func startTerminal(v bool) error {
	log.Println("skewer start.")
	verbose = v
	return nil
}

func endTerminal(bin string) {
	cleanup(bin)
}

func printVerbose(args ...interface{}) {
	if verbose {
		log.Println(args...)
	}
}

func setCommandPipe(cmd *exec.Cmd, test bool) error {

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return xerrors.Errorf("StdoutPipe error: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return xerrors.Errorf("StdoutPipe error: %w", err)
	}

	setPipe(36, stdout, os.Stdout, test)
	setPipe(31, stderr, os.Stderr, test)
	return nil
}

func setPipe(c int, r io.ReadCloser, w io.Writer, test bool) {
	go func() {
		s := bufio.NewScanner(r)

		for s.Scan() {

			txt := s.Text()
			co := c

			if test {
				idx := strings.Index(txt, "PASS")
				ok := strings.Index(txt, "ok")
				if idx >= 0 && idx < 7 || ok == 0 {
					co = 32
				} else {
					idx = strings.Index(txt, "FAIL")
					if idx >= 0 && idx < 7 {
						co = 31
					}
				}
			}

			fmt.Fprintf(w, "\x1b[%dm| %s\x1b[0m\n", co, txt)
		}
	}()
	return
}
