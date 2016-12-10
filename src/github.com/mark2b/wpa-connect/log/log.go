package log

import (
	"github.com/op/go-logging"
	"os"
)

var (
	Log = logging.MustGetLogger("")
)

func SetSilentMode() {
	initialize(defaultModeFormatter(), logging.WARNING)
}

func SetInfoMode() {
	initialize(defaultModeFormatter(), logging.INFO)
}

func SetVerboseMode() {
	initialize(defaultModeFormatter(), logging.DEBUG)
}

func SetDebugMode() {
	initialize(debugModeFormatter(), logging.DEBUG)
}

func initialize(formatter logging.Formatter, level logging.Level) {
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, formatter)
	logging.SetBackend(backendFormatter)
	logging.SetLevel(level, "")
}

func init() {
	SetInfoMode()
}

func debugModeFormatter() logging.Formatter {
	return logging.MustStringFormatter(
		`%{color}%{level:.1s} %{module} %{time:15:04:05.000} %{shortfile} %{message}%{color:reset}`,
	)
}

func defaultModeFormatter() logging.Formatter {
	return logging.MustStringFormatter(
		`%{level:.1s} %{module} %{time:15:04:05.000} %{message}`,
	)
}