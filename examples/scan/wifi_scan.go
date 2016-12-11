package main

import (
	wifi "github.com/mark2b/wpa-connect"
)

func main2() {
	if bssList, err := wifi.ScanManager.Scan(); err == nil {
		for _, bss := range bssList {
			print(bss.SSID, bss.Signal, bss.KeyMgmt)
		}
	}
}
