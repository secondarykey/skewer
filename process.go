package skewer

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

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

	//TODO retry

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

func checkConnection(port int) (chan bool, error) {

	ch := make(chan bool)

	//TODO ???

	go func() {
		var conn net.Conn
		var err error
		for range time.Tick(100 * time.Millisecond) {
			conn, err = net.Dial("tcp", fmt.Sprintf(":%d", port))
			if err == nil {
				break
			}
			log.Println(err)
		}
		defer conn.Close()
		ch <- true
		close(ch)
	}()

	return ch, nil
}

func cleanup(bin string) {

	err := kill()
	if err != nil {
		log.Println(err)
	}

	// TODO process kill waiting
	time.Sleep(1 * time.Second)

	if _, err := os.Stat(bin); err == nil {
		err = os.Remove(bin)
		if err != nil {
			log.Println(err)
		}
	}
}

func startServer(bin string, port int, args []string, ch chan error) {

	log.Println("Start Build and Launch.")

	setStatus(BuildStatus)

	err := build(bin, args)
	if err != nil {
		log.Println("Build Error")
		setStatus(BuildErrorStatus)
		ch <- err
		return
	}

	setStatus(StartupStatus)

	// run process
	err = run(bin)
	if err != nil {
		log.Println("Process Run Error")
		setStatus(StartupErrorStatus)
		ch <- err
		return
	}

	if port == 0 {
		log.Println("Complete Build and Launch")
		setStatus(OKStatus)
	} else {
		go func() {
			check, err := checkConnection(port)

			if err != nil {
				log.Println("HTTP Connection Error")
				setStatus(StartupErrorStatus)
				//TODO ?
				ch <- err
				return
			}

			v := <-check
			if v {
				log.Println("Complete Build and Launch")
				setStatus(OKStatus)
			} else {
				log.Println("HTTP Server Launch Error")
				setStatus(StartupErrorStatus)
			}
		}()
	}
}
