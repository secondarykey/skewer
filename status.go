package skewer

import (
	"sync"
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
var statusMutex sync.Mutex

func getStatus() Status {
	statusMutex.Lock()
	defer statusMutex.Unlock()
	return status
}

func setStatus(s Status) {
	statusMutex.Lock()
	defer statusMutex.Unlock()
	if status == FatalStatus {
		return
	}
	status = s
	printVerbose("Status:", s)
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
