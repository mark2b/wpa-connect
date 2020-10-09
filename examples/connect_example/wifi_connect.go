package main

import (
	"fmt"
	"os"

	wifi "wpa-connect"
	"time"
)

func main() {

	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("Insufficient arguments")
		return
	}
	ssid := args[0]
	password := args[1]
	wifi.SetDebugMode()
	if conn, err := wifi.ConnectManager.Connect(ssid, password, time.Second * 60); err == nil {
		fmt.Println("Connected", conn.NetInterface, conn.SSID, conn.IP4.String(), conn.IP6.String())
	} else {
		fmt.Println(err)
	}
}
