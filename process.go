package skewer

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"golang.org/x/xerrors"
)

var current *os.Process
var processMutex sync.Mutex

func run(name string) error {

	wd, err := os.Getwd()
	if err != nil {
		return xerrors.Errorf("work directory get error: %w", err)
	}

	cmd := exec.Command(filepath.Join(wd, name))

	err = setCommandPipe(cmd)
	if err != nil {
		return xerrors.Errorf("setCommandPipe() error: %w", err)
	}

	err = cmd.Start()
	if err != nil {
		return xerrors.Errorf("http server start error: %w", err)
	}

	processMutex.Lock()
	defer processMutex.Unlock()
	current = cmd.Process

	return nil
}

func kill() error {
	//TODO リトライする

	processMutex.Lock()
	defer processMutex.Unlock()
	if current != nil {
		err := current.Kill()
		if err != nil {
			return xerrors.Errorf("current process kill error: %w", err)
		}
	}

	return nil
}
