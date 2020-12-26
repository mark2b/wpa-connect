package main

import (
	wifi "github.com/mark2b/wpa-connect"
)

func main() {
	if bssList, err := wifi.ScanManager.Scan(); err == nil {
		for _, bss := range bssList {
			print(bss.SSID, bss.Signal, bss.KeyMgmt)
		}
	}
}
