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

	if current != nil {
		err := Kill()
		if err != nil {
			return err
		}
	}

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

	//TODO 起動中のエラーを取りたい
	current = cmd.Process

	//TODO 起動監視

	return nil
}

func Kill() error {
	if current != nil {
		err := current.Kill()
		if err != nil {
			return xerrors.Errorf("current process kill error: %w", err)
		}
	}
	return nil
}
