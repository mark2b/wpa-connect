package log

import (
	"github.com/op/go-logging"
	"os"
)

var (
	Log = logging.MustGetLogger("")
)

func SetVerboseMode() {
	logging.SetLevel(logging.DEBUG, "")
}

func SetDebugMode() {
	format := logging.MustStringFormatter(
		`%{color}%{level:.1s} %{module} %{time:15:04:05.000} %{shortfile} %{message}%{color:reset}`,
	)
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
	logging.SetLevel(logging.DEBUG, "")
}

func Start() {
	logging.SetLevel(logging.INFO, "")
}

func init() {
	format := logging.MustStringFormatter(
		`%{level:.1s} %{module} %{time:15:04:05.000} %{message}`,
	)
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
	logging.SetLevel(logging.INFO, "")

	SetDebugMode()
	SetVerboseMode()
}
