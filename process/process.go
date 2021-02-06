package process

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/xerrors"
)

var current *os.Process

func Run(name string) error {

	wd, err := os.Getwd()
	if err != nil {
		return xerrors.Errorf("work directory get error: %w", err)
	}

	cmd := exec.Command(filepath.Join(wd, name))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println("HTTP Server Start")

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
