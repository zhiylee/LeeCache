package peers

const (
	StatusStopped = 0
	StatusOK      = 200
	StatusError   = 500
)

var StatusText = map[int]string{
	StatusStopped: "Stopped",
	StatusOK:      "Running",
	StatusError:   "Error",
}
