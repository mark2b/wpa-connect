package wpa_cli

import (
	"fmt"
	"os/exec"
)

type WPACli struct {
	NetInterface string
}

func (self *WPACli) SaveConfig() (e error) {
	cmd := exec.Command("wpa_cli", fmt.Sprintf("-i%s", self.NetInterface), "save_config")

	return cmd.Start()
}
