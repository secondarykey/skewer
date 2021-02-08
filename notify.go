package skewer

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/xerrors"
)

const ModFile = "go.mod"

func notifyMonitoring(args []string, patterns []string, ch chan error) {

	mod := searchPath(args)
	if mod == "" {
		setStatus(FatalStatus)
		ch <- fmt.Errorf("Not found go.mod file.")
		return
	}

	printVerbose("Specification:", mod)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		setStatus(FatalStatus)
		ch <- xerrors.Errorf("fsnotify.NewWatcher() error: %w", err)
		return
	}
	defer watcher.Close()

	err = registerWatcher(watcher, mod)
	if err != nil {
		setStatus(FatalStatus)
		ch <- xerrors.Errorf("registerWatcher() error: %w", err)
		return
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				ch <- fmt.Errorf("watcher event not OK")
				continue
			}

			if !ignoreFile(event.Name, patterns) {
				if event.Op&fsnotify.Write == fsnotify.Write {
					//log.Println("modified file:", event.Name)
				}

				s := getStatus()
				if s.reboot() {
					printVerbose(event.Name)
					log.Println("Waiting for reboot.")
					setStatus(WaitingForRebootStatus)
				}
			}
		case err, _ := <-watcher.Errors:
			if err != nil {
				ch <- xerrors.Errorf("watcher errors: %w", err)
			}
		}
	}

	return
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

func searchPath(patterns []string) string {

	for _, pattern := range patterns {

		files, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}

		for _, file := range files {
			abs, err := filepath.Abs(file)
			if err != nil {
				continue
			}

			dir := filepath.Dir(abs)

			info, err := os.Stat(abs)
			if err != nil {
				continue
			}
			if info.IsDir() {
				dir = abs
			}
			return searchModFile(dir)
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
