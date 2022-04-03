package main

import (
	"fmt"
	"os"
	"time"

	wifi "github.com/mark2b/wpa-connect"
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

	conn, err := wifi.ConnectManager.Connect(ssid, password, time.Second*60)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected", conn.NetInterface, conn.SSID, conn.IP4.String(), conn.IP6.String())
}
