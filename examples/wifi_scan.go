package main

import (
	wifi "github.com/mark2b/wpaconnect"
	"github.com/mark2b/wpaconnect/log"
)

func main() {
	if bssList, err := wifi.ScanManager.Scan(); err == nil {
		for _, bss := range bssList {
			log.Log.Info(bss.SSID, bss.Signal, bss.KeyMgmt)
		}
	}
}
