package main

import (
	"fmt"
	"os"

	wifi "github.com/mark2b/wpa-connect"
	"github.com/mark2b/wpa-connect/log"
)

func main() {
	log.SetDebugMode()
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("Insufficient arguments")
		return
	}
	ssid := args[0]
	password := args[1]
	if err := wifi.ConnectManager.Connect(ssid, password); err == nil {
		fmt.Println("Connected")
	} else {
		fmt.Println(err)
	}
}
