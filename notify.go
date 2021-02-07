package skewer

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/secondarykey/skewer/terminal"
	"golang.org/x/xerrors"
)

func monitoring(args []string, patterns []string) error {

	mod := searchPath(args)
	if mod == "" {
		return fmt.Errorf("not found go.mod file.")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return xerrors.Errorf("fsnotify.NewWatcher() error: %w", err)
	}
	defer watcher.Close()

	err = registerWatcher(watcher, mod)
	if err != nil {
		return xerrors.Errorf("registerWatcher() error: %w", err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return fmt.Errorf("not ok")
			}

			if !ignoreFile(event.Name, patterns) {
				if event.Op&fsnotify.Write == fsnotify.Write {
					//log.Println("modified file:", event.Name)
				}

				s := getStatus()
				if s.reboot() {
					terminal.Verbose(event.Name)
					log.Println("Waiting for reboot.")
					setStatus(WaitingForRebootStatus)
				}
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

func ignoreFile(path string, patterns []string) bool {

	name := filepath.Base(path)
	for _, pattern := range patterns {
		if match, err := filepath.Match(pattern, name); match {
			return true
		} else if err != nil {
			log.Println(err)
		}
	}
	return false
}

func registerWatcher(w *fsnotify.Watcher, path string) error {
	err := w.Add(path)
	if err != nil {
		return xerrors.Errorf("watcher.Add() error: %w", err)
	}
	entry, err := os.ReadDir(path)
	if err != nil {
		return xerrors.Errorf("os.ReadDir() error: %w", err)
	}

	for _, info := range entry {
		if info.IsDir() {
			err := registerWatcher(w, filepath.Join(path, info.Name()))
			if err != nil {
				return xerrors.Errorf("registerWatcher() error: %w", err)
			}
		}
	}
	return nil
}

const ModFile = "go.mod"

func searchPath(files []string) string {

	for _, elm := range files {

		abs, err := filepath.Abs(elm)
		if err == nil {
			dir := filepath.Dir(abs)
			return searchModFile(dir)
		} else {
			//TODO Glob
		}
	}

	return ""
}

func searchModFile(dir string) string {

	entry, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	for _, elm := range entry {
		if elm.Name() == ModFile {
			return dir
		}
	}

	rtn := filepath.Dir(dir)
	if rtn == dir {
		return ""
	}
	return searchModFile(rtn)
}
