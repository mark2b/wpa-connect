package main

import (
	"fmt"

	wifi "github.com/mark2b/wpa-connect"
)

func main() {
	bssList, err := wifi.ScanManager.Scan()
	if err != nil {
		panic(err)
	}

	for _, bss := range bssList {
		fmt.Println(bss.SSID, bss.Signal, bss.KeyMgmt)
	}
}
