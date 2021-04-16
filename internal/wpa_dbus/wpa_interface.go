package wpa_dbus

import (
	"fmt"

	"github.com/godbus/dbus"
	"github.com/mark2b/wpa-connect/internal/log"
)

type InterfaceWPA struct {
	WPA              *WPA
	Object           dbus.BusObject
	Networks         []NetworkWPA
	BSSs             []BSSWPA
	State            string
	Scanning         bool
	Ifname           string
	CurrentBSS       *BSSWPA
	TempBSS          *BSSWPA
	CurrentNetwork   *NetworkWPA
	NewNetwork       *NetworkWPA
	ScanInterval     int32
	DisconnectReason int32
	SignalChannel    chan *dbus.Signal
	Error            error
}

func (self *InterfaceWPA) ReadNetworksList() *InterfaceWPA {
	if self.Error == nil {
		if networks, err := self.WPA.get("fi.w1.wpa_supplicant1.Interface.Networks", self.Object); err == nil {
			newNetworks := []NetworkWPA{}
			for _, networkObjectPath := range networks.([]dbus.ObjectPath) {
				network := NetworkWPA{Interface: self, Object: self.WPA.Connection.Object("fi.w1.wpa_supplicant1", networkObjectPath)}
				newNetworks = append(newNetworks, network)
			}
			self.Networks = newNetworks
		} else {
			self.Error = err
		}
	}
	return self
}

func (self *InterfaceWPA) MakeTempBSS() *InterfaceWPA {
	self.TempBSS = &BSSWPA{Interface: self, Object: self.WPA.Connection.Object("fi.w1.wpa_supplicant1", "/")}
	return self
}

func (self *InterfaceWPA) ReadBSSList() *InterfaceWPA {
	if self.Error == nil {
		if bsss, err := self.WPA.get("fi.w1.wpa_supplicant1.Interface.BSSs", self.Object); err == nil {
			newBSSs := []BSSWPA{}
			for _, bssObjectPath := range bsss.([]dbus.ObjectPath) {
				bss := BSSWPA{Interface: self, Object: self.WPA.Connection.Object("fi.w1.wpa_supplicant1", bssObjectPath)}
				newBSSs = append(newBSSs, bss)
			}
			self.BSSs = newBSSs
		} else {
			self.Error = err
		}
	}
	return self
}

func (self *InterfaceWPA) Scan() *InterfaceWPA {
	if self.Error == nil {
		args := make(map[string]dbus.Variant, 0)
		args["Type"] = dbus.MakeVariant("passive")
		if call := self.Object.Call("fi.w1.wpa_supplicant1.Interface.Scan", 0, args); call.Err == nil {
		} else {
			self.Error = call.Err
		}
	}
	return self
}

func (self *InterfaceWPA) Disconnect() *InterfaceWPA {
	if self.Error == nil {
		if call := self.Object.Call("fi.w1.wpa_supplicant1.Interface.Disconnect", 0); call.Err == nil {
		} else {
			self.Error = call.Err
		}
	}
	return self
}

func (self *InterfaceWPA) Reassociate() *InterfaceWPA {
	if self.Error == nil {
		if call := self.Object.Call("fi.w1.wpa_supplicant1.Interface.Reassociate", 0); call.Err == nil {
		} else {
			self.Error = call.Err
		}
	}
	return self
}

func (self *InterfaceWPA) Reattach() *InterfaceWPA {
	if self.Error == nil {
		if call := self.Object.Call("fi.w1.wpa_supplicant1.Interface.Reattach", 0); call.Err == nil {
		} else {
			self.Error = call.Err
		}
	}
	return self
}

func (self *InterfaceWPA) Reconnect() *InterfaceWPA {
	if self.Error == nil {
		if call := self.Object.Call("fi.w1.wpa_supplicant1.Interface.Reconnect", 0); call.Err == nil {
		} else {
			self.Error = call.Err
		}
	}
	return self
}

func (self *InterfaceWPA) RemoveAllNetworks() *InterfaceWPA {
	if self.Error == nil {
		if call := self.Object.Call("fi.w1.wpa_supplicant1.Interface.RemoveAllNetworks", 0); call.Err == nil {
		} else {
			self.Error = call.Err
		}
	}
	return self
}

func (self *InterfaceWPA) AddNetwork(args map[string]dbus.Variant) *InterfaceWPA {
	if self.Error == nil {
		if call := self.Object.Call("fi.w1.wpa_supplicant1.Interface.AddNetwork", 0, args); call.Err == nil {
			if len(call.Body) > 0 {
				networkObjectPath := call.Body[0].(dbus.ObjectPath)
				self.NewNetwork = &NetworkWPA{Interface: self, Object: self.WPA.Connection.Object("fi.w1.wpa_supplicant1", networkObjectPath)}
			}
		} else {
			self.Error = call.Err
		}
	}
	return self
}

func (self *InterfaceWPA) ReadState() *InterfaceWPA {
	if self.Error == nil {
		if value, err := self.WPA.get("fi.w1.wpa_supplicant1.Interface.State", self.Object); err == nil {
			self.State = value.(string)
		} else {
			self.Error = err
		}
	}
	return self
}

func (self *InterfaceWPA) ReadScanning() *InterfaceWPA {
	if self.Error == nil {
		if value, err := self.WPA.get("fi.w1.wpa_supplicant1.Interface.Scanning", self.Object); err == nil {
			self.Scanning = value.(bool)
		} else {
			self.Error = err
		}
	}
	return self
}

func (self *InterfaceWPA) ReadIfname() *InterfaceWPA {
	if self.Error == nil {
		if value, err := self.WPA.get("fi.w1.wpa_supplicant1.Interface.Ifname", self.Object); err == nil {
			self.Ifname = value.(string)
		} else {
			self.Error = err
		}
	}
	return self
}

func (self *InterfaceWPA) ReadScanInterval() *InterfaceWPA {
	if self.Error == nil {
		if value, err := self.WPA.get("fi.w1.wpa_supplicant1.Interface.ScanInterval", self.Object); err == nil {
			self.ScanInterval = value.(int32)
		} else {
			self.Error = err
		}
	}
	return self
}

func (self *InterfaceWPA) ReadDisconnectReason() *InterfaceWPA {
	if self.Error == nil {
		if value, err := self.WPA.get("fi.w1.wpa_supplicant1.Interface.DisconnectReason", self.Object); err == nil {
			self.DisconnectReason = value.(int32)
		} else {
			self.Error = err
		}
	}
	return self
}

func (self *InterfaceWPA) AddSignalsObserver() *InterfaceWPA {
	log.Log.Debug("AddSignalsObserver.Interface")
	match := fmt.Sprintf("type='signal',interface='fi.w1.wpa_supplicant1.Interface',path='%s'", self.Object.Path())
	if call := self.WPA.Connection.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, match); call.Err == nil {
	} else {
		self.Error = call.Err
	}
	return self
}

func (self *InterfaceWPA) RemoveSignalsObserver() *InterfaceWPA {
	log.Log.Debug("RemoveSignalsObserver.Interface")
	match := fmt.Sprintf("type='signal',interface='fi.w1.wpa_supplicant1.Interface',path='%s'", self.Object.Path())
	if call := self.WPA.Connection.BusObject().Call("org.freedesktop.DBus.RemoveMatch", 0, match); call.Err == nil {
	} else {
		self.Error = call.Err
	}
	return self
}

func (self *InterfaceWPA) ReadCurrentBSS() *InterfaceWPA {
	if self.Error == nil {
		if value, err := self.WPA.get("fi.w1.wpa_supplicant1.Interface.CurrentBSS", self.Object); err == nil {
			bssObjectPath := value.(dbus.ObjectPath)
			self.CurrentBSS = &BSSWPA{Interface: self, Object: self.WPA.Connection.Object("fi.w1.wpa_supplicant1", bssObjectPath)}
		} else {
			self.Error = err
		}
	}
	return self
}

func (self *InterfaceWPA) ReadCurrentNetwork() *InterfaceWPA {
	if self.Error == nil {
		if network, err := self.WPA.get("fi.w1.wpa_supplicant1.Interface.CurrentNetwork", self.Object); err == nil {
			networkObjectPath := network.(dbus.ObjectPath)
			self.CurrentNetwork = &NetworkWPA{Interface: self, Object: self.WPA.Connection.Object("fi.w1.wpa_supplicant1", networkObjectPath)}
		} else {
			self.Error = err
		}
	}
	return self
}
