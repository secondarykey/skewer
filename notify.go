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

func searchWatchingPath(args []string) ([]string, error) {

	mod := searchPath(args)
	if mod == "" {
		return nil, xerrors.Errorf("NotFound go.mod file -> %v", args)
	}
	printVerbose("Specification:", mod)

	paths, err := getDirectories(mod)
	if err != nil {
		return nil, xerrors.Errorf("getDirectories() error: %w", err)
	}

	return paths, nil
}

func notifyMonitoring(paths []string, patterns []string, ch chan error) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		setStatus(FatalStatus)
		ch <- xerrors.Errorf("fsnotify.NewWatcher() error: %w", err)
		return
	}
	defer watcher.Close()

	for _, path := range paths {
		err = watcher.Add(path)
		if err != nil {
			setStatus(FatalStatus)
			ch <- xerrors.Errorf("watcher Add() error: %w", err)
			return
		}
	}

	for {
		select {
		case event, ok := <-watcher.Events:

			if !ok {
				ch <- fmt.Errorf("watcher event not OK")
				continue
			}

			//TODO new directory
			//log.Println("event file:", event.Name)

			if !ignoreFile(event.Name, patterns) {

				if event.Op&fsnotify.Write == fsnotify.Write {
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

func getDirectories(path string) ([]string, error) {

	paths := make([]string, 0, 100)
	paths = append(paths, path)

	entry, err := os.ReadDir(path)
	if err != nil {
		return nil, xerrors.Errorf("os.ReadDir() error: %w", err)
	}

	for _, info := range entry {
		if info.IsDir() {
			work, err := getDirectories(filepath.Join(path, info.Name()))
			if err != nil {
				return nil, xerrors.Errorf("getDirectories() error: %w", err)
			}
			paths = append(paths, work...)
		}
	}

	return paths, nil
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
