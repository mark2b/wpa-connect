package main

import (
	wifi "github.com/mark2b/wpa-connect"
	"github.com/mark2b/wpa-connect/log"
)

func main() {
	if bssList, err := wifi.ScanManager.Scan(); err == nil {
		for _, bss := range bssList {
			log.Log.Info(bss.SSID, bss.Signal, bss.KeyMgmt)
		}
	}
}
