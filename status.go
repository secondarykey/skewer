package skewer

type Status int

const (
	ReadyStatus Status = iota
	FatalStatus
	BuildStatus
	BuildErrorStatus
	StartupStatus
	StartupErrorStatus
	OKStatus
)

var status Status = ReadyStatus

func getStatus() Status {
	return status
}

func setStatus(s Status) {
	status = s
}
