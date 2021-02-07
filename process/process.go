package process

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/secondarykey/skewer/terminal"
	"golang.org/x/xerrors"
)

var current *os.Process

func Run(name string) error {

	wd, err := os.Getwd()
	if err != nil {
		return xerrors.Errorf("work directory get error: %w", err)
	}

	cmd := exec.Command(filepath.Join(wd, name))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return xerrors.Errorf("StdoutPipe error: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return xerrors.Errorf("StdoutPipe error: %w", err)
	}
	terminal.SetPipe(stdout, stderr)

	err = cmd.Start()
	if err != nil {
		return xerrors.Errorf("http server start error: %w", err)
	}
	current = cmd.Process

	return nil
}

func Kill() error {
	//TODO リトライする
	if current != nil {
		err := current.Kill()
		if err != nil {
			return xerrors.Errorf("current process kill error: %w", err)
		}
	}
	return nil
}
