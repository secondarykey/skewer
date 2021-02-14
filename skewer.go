package skewer

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/secondarykey/skewer/config"

	"golang.org/x/xerrors"
)

func Patrol(opts ...config.Option) error {

	if !checkGo() {
		return xerrors.Errorf(`requires "go" to skewer run.`)
	}

	err := config.Set(opts)
	if err != nil {
		return xerrors.Errorf("config.Set() error: %w", err)
	}

	conf := config.Get()
	bin := conf.Bin

	ch := make(chan error)
	done := make(chan error)

	startTerminal(conf.Verbose)
	defer endTerminal(bin)

	// error signal
	go func() {
		for {
			select {
			case err := <-ch:
				if err != nil {
					msg := fmt.Sprintf("%+v", err)
					printVerbose(msg)
					if getStatus() == FatalStatus {
						done <- err
						return
					}
				}
			default:
			}
		}
	}()

	// fsnotify
	go func() {
		notifyMonitoring(conf.Args, conf.IgnoreFiles, ch)
	}()

	// Ctl + c Signal
	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		<-quit
		done <- nil
	}()

	// rebuild
	go func() {
		setStatus(WaitingForRebootStatus)
		rebuildMonitor(5, ch)
	}()

	// wait
	return <-done
}
