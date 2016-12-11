package wpaconnect

import "github.com/mark2b/wpa-connect/internal/log"

func SetSilentMode() {
	log.SetSilentMode()
}

func SetInfoMode() {
	log.SetInfoMode()
}

func SetVerboseMode() {
	log.SetVerboseMode()
}

func SetDebugMode() {
	log.SetDebugMode()
}
