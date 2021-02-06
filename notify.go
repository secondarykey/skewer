package skewer

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/xerrors"
)

func monitoring() error {

	path := "C:\\Users\\secon\\GoApp\\section"
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return xerrors.Errorf("fsnotify.NewWatcher() error: %w", err)
	}
	defer watcher.Close()

	err = watcher.Add(path)
	if err != nil {
		return xerrors.Errorf("watcher.Add() error: %w", err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return fmt.Errorf("not ok")
			}

			log.Println("event:", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("modified file:", event.Name)
			}

			if getStatus() == OKStatus {
				setStatus(WaitingForRebootStatus)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return fmt.Errorf("not ok")
			}
			log.Println("error:", err)
		}
	}

	return nil
}
