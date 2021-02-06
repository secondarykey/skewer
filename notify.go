package skewer

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/xerrors"
)

func monitoring(args []string) error {

	mod := searchPath(args)
	if mod == "" {
		return fmt.Errorf("not found go.mod file.")
	}

	fmt.Println(mod)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return xerrors.Errorf("fsnotify.NewWatcher() error: %w", err)
	}
	defer watcher.Close()

	err = watcher.Add(mod)
	if err != nil {
		return xerrors.Errorf("watcher.Add() error: %w", err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return fmt.Errorf("not ok")
			}

			//TODO ignore

			//log.Println("event:", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				//log.Println("modified file:", event.Name)
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

const ModFile = "go.mod"

func searchPath(files []string) string {

	for _, elm := range files {
		abs, err := filepath.Abs(elm)
		if err == nil {
			if err == nil {
				dir := searchModFile(abs)
				if dir != "" {
					return dir
				}
			} else {
				//TODO Glob
			}
		}

	}

	return ""
}

func searchModFile(file string) string {
	info, err := os.Stat(file)
	if err != nil {
		return ""
	}
	if err == nil {
		name := info.Name()
		if info.IsDir() {
			infos, err := os.ReadDir(name)
			if err != nil {
				return ""
			}
			for _, elm := range infos {
				dir := searchModFile(elm.Name())
				if dir != "" {
					return dir
				}
			}
		} else {
			dir := filepath.Dir(name)
			base := filepath.Base(name)
			if base == ModFile {
				return dir
			}

			rtn := searchModFile(dir)
			if rtn != "" {
				return rtn
			}
		}
	}
	return ""
}
