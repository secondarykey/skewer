package skewer

import (
	"time"

	"github.com/secondarykey/skewer/config"
	"github.com/secondarykey/skewer/terminal"
)

type Status int

const (
	ReadyStatus Status = iota
	FatalStatus
	BuildStatus
	BuildErrorStatus
	StartupStatus
	StartupErrorStatus
	WaitingForRebootStatus
	OKStatus
)

var status Status = ReadyStatus

func getStatus() Status {
	return status
}

func setStatus(s Status) {
	status = s
	terminal.Verbose("Status:", s)
}

func (s Status) reboot() bool {
	if s == OKStatus || s == BuildErrorStatus || s == StartupErrorStatus {
		return true
	}
	return false
}

func (s Status) canBuild() bool {
	if s == WaitingForRebootStatus {
		return true
	}
	return false
}

func (s Status) String() string {
	switch s {
	case ReadyStatus:
		return "skewer process start."
	case FatalStatus:
		return "Unable to continue process."
	case BuildStatus:
		return "Building the target."
	case BuildErrorStatus:
		return "Build error."
	case StartupStatus:
		return "HTTP Server start up."
	case StartupErrorStatus:
		return "HTTP Server start up error."
	case WaitingForRebootStatus:
		return "HTTP Server waiting for restart(rebuild)."
	case OKStatus:
		return "HTTP Server Accessable."
	}
	return "NotFound Status"
}

func rebuildMonitor(s int, ch chan error) error {

	conf := config.Get()
	bin := conf.Bin

	//初回起動を行う
	cleanup(bin)
	go startServer(bin, conf.AppPort, conf.Args, ch)

	d := time.Duration(s) * time.Second
	for range time.Tick(d) {
		status = getStatus()
		//ビルド待ちだった場合
		if status.canBuild() {
			cleanup(bin)
			go startServer(bin, conf.AppPort, conf.Args, ch)
		}
	}

	return nil
}
