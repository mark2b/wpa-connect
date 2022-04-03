package wpa_dbus

import (
	"fmt"

	"github.com/godbus/dbus"
	"github.com/mark2b/wpa-connect/internal/log"
)

type NetworkWPA struct {
	Interface     *InterfaceWPA
	Object        dbus.BusObject
	SSID          string
	KeyMgmt       string
	SignalChannel chan *dbus.Signal
	Error         error
}

func (self *NetworkWPA) ReadProperties() *NetworkWPA {
	log.Log.Debug("ReadProperties")
	if self.Error != nil {
		return self
	}

	properties, err := self.Interface.WPA.get("fi.w1.wpa_supplicant1.Network.Properties", self.Object)
	if err != nil {
		self.Error = err
		return self
	}

	for key, value := range properties.(map[string]dbus.Variant) {
		switch key {
		case "ssid":
			self.SSID = value.Value().(string)
		case "key_mgmt":
			self.KeyMgmt = value.Value().(string)
		}
	}

	return self
}

func (self *NetworkWPA) Select() *NetworkWPA {
	log.Log.Debug("Select")
	if self.Error != nil {
		return self
	}

	call := self.Interface.Object.Call("fi.w1.wpa_supplicant1.Interface.SelectNetwork", 0, self.Object.Path())
	self.Error = call.Err

	return self
}

func (self *NetworkWPA) AddSignalsObserver() *NetworkWPA {
	log.Log.Debug("AddSignalsObserver.Network")
	match := fmt.Sprintf("type='signal',interface='fi.w1.wpa_supplicant1.Network',path='%s'", self.Object.Path())
	call := self.Interface.WPA.Connection.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, match)
	self.Error = call.Err

	return self
}

func (self *NetworkWPA) RemoveSignalsObserver() *NetworkWPA {
	log.Log.Debug("RemoveSignalsObserver.Network")
	match := fmt.Sprintf("type='signal',interface='fi.w1.wpa_supplicant1.Network',path='%s'", self.Object.Path())
	call := self.Interface.WPA.Connection.BusObject().Call("org.freedesktop.DBus.RemoveMatch", 0, match)
	self.Error = call.Err

	return self
}
