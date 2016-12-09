package wpa_cli

import (
	"fmt"

	"github.com/ThomasRooney/gexpect"
)

type WPACli struct {
	NetInterface string
}

func (self *WPACli) SaveConfig() (e error) {
	cmd := fmt.Sprintf("wpa_cli -i%s save_config", self.NetInterface)
	if _, err := gexpect.Spawn(cmd); err == nil {
	} else {
		e = err
	}
	return
}
