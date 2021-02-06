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
		return "Proxy Server Accessable."
	}
	return "NotFound Status"
}

func rebuildMonitor(s int) error {

	conf := config.Get()
	bin := conf.Bin

	d := time.Duration(s) * time.Second
	for range time.Tick(d) {
		status := getStatus()
		//ビルド待ちだった場合
		if status == WaitingForRebootStatus {
			cleanup(bin)
			go startServer(bin, conf.AppPort, conf.Args)
		}
	}

	return nil
}
