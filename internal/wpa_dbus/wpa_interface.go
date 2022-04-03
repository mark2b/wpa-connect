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
	if self.Error != nil {
		return self
	}

	networks, err := self.WPA.get("fi.w1.wpa_supplicant1.Interface.Networks", self.Object)
	if err != nil {
		self.Error = err
		return self
	}

	newNetworks := []NetworkWPA{}
	for _, networkObjectPath := range networks.([]dbus.ObjectPath) {
		network := NetworkWPA{Interface: self, Object: self.WPA.Connection.Object("fi.w1.wpa_supplicant1", networkObjectPath)}
		newNetworks = append(newNetworks, network)
	}
	self.Networks = newNetworks

	return self
}

func (self *InterfaceWPA) MakeTempBSS() *InterfaceWPA {
	self.TempBSS = &BSSWPA{
		Interface: self,
		Object:    self.WPA.Connection.Object("fi.w1.wpa_supplicant1", "/"),
	}
	return self
}

func (self *InterfaceWPA) ReadBSSList() *InterfaceWPA {
	if self.Error != nil {
		return self
	}

	bsss, err := self.WPA.get("fi.w1.wpa_supplicant1.Interface.BSSs", self.Object)
	if err != nil {
		self.Error = err
		return self
	}

	newBSSs := []BSSWPA{}
	for _, bssObjectPath := range bsss.([]dbus.ObjectPath) {
		bss := BSSWPA{Interface: self, Object: self.WPA.Connection.Object("fi.w1.wpa_supplicant1", bssObjectPath)}
		newBSSs = append(newBSSs, bss)
	}
	self.BSSs = newBSSs

	return self
}

func (self *InterfaceWPA) Scan() *InterfaceWPA {
	if self.Error != nil {
		return self
	}

	args := make(map[string]dbus.Variant, 0)
	args["Type"] = dbus.MakeVariant("passive")
	call := self.Object.Call("fi.w1.wpa_supplicant1.Interface.Scan", 0, args)
	self.Error = call.Err

	return self
}

func (self *InterfaceWPA) Disconnect() *InterfaceWPA {
	if self.Error != nil {
		return self
	}

	call := self.Object.Call("fi.w1.wpa_supplicant1.Interface.Disconnect", 0)
	self.Error = call.Err

	return self
}

func (self *InterfaceWPA) Reassociate() *InterfaceWPA {
	if self.Error != nil {
		return self
	}

	call := self.Object.Call("fi.w1.wpa_supplicant1.Interface.Reassociate", 0)
	self.Error = call.Err

	return self
}

func (self *InterfaceWPA) Reattach() *InterfaceWPA {
	if self.Error != nil {
		return self
	}

	call := self.Object.Call("fi.w1.wpa_supplicant1.Interface.Reattach", 0)
	self.Error = call.Err

	return self
}

func (self *InterfaceWPA) Reconnect() *InterfaceWPA {
	if self.Error != nil {
		return self
	}

	call := self.Object.Call("fi.w1.wpa_supplicant1.Interface.Reconnect", 0)
	self.Error = call.Err

	return self
}

func (self *InterfaceWPA) RemoveAllNetworks() *InterfaceWPA {
	if self.Error != nil {
		return self
	}

	call := self.Object.Call("fi.w1.wpa_supplicant1.Interface.RemoveAllNetworks", 0)
	self.Error = call.Err

	return self
}

func (self *InterfaceWPA) AddNetwork(args map[string]dbus.Variant) *InterfaceWPA {
	if self.Error != nil {
		return self
	}

	call := self.Object.Call("fi.w1.wpa_supplicant1.Interface.AddNetwork", 0, args)
	if call.Err != nil {
		self.Error = call.Err
		return self
	}

	if len(call.Body) > 0 {
		networkObjectPath := call.Body[0].(dbus.ObjectPath)
		self.NewNetwork = &NetworkWPA{Interface: self, Object: self.WPA.Connection.Object("fi.w1.wpa_supplicant1", networkObjectPath)}
	}

	return self
}

func (self *InterfaceWPA) ReadState() *InterfaceWPA {
	if self.Error != nil {
		return self
	}

	value, err := self.WPA.get("fi.w1.wpa_supplicant1.Interface.State", self.Object)
	if err != nil {
		self.Error = err
		return self
	}
	self.State = value.(string)

	return self
}

func (self *InterfaceWPA) ReadScanning() *InterfaceWPA {
	if self.Error != nil {
		return self
	}

	value, err := self.WPA.get("fi.w1.wpa_supplicant1.Interface.Scanning", self.Object)
	if err != nil {
		self.Error = err
		return self
	}
	self.Scanning = value.(bool)

	return self
}

func (self *InterfaceWPA) ReadIfname() *InterfaceWPA {
	if self.Error != nil {
		return self
	}

	value, err := self.WPA.get("fi.w1.wpa_supplicant1.Interface.Ifname", self.Object)
	if err != nil {
		self.Error = err
		return self
	}
	self.Ifname = value.(string)

	return self
}

func (self *InterfaceWPA) ReadScanInterval() *InterfaceWPA {
	if self.Error != nil {
		return self
	}

	value, err := self.WPA.get("fi.w1.wpa_supplicant1.Interface.ScanInterval", self.Object)
	if err != nil {
		self.Error = err
		return self
	}
	self.ScanInterval = value.(int32)

	return self
}

func (self *InterfaceWPA) ReadDisconnectReason() *InterfaceWPA {
	if self.Error != nil {
		return self
	}

	value, err := self.WPA.get("fi.w1.wpa_supplicant1.Interface.DisconnectReason", self.Object)
	if err != nil {
		self.Error = err
		return self
	}
	self.DisconnectReason = value.(int32)

	return self
}

func (self *InterfaceWPA) AddSignalsObserver() *InterfaceWPA {
	log.Log.Debug("AddSignalsObserver.Interface")

	match := fmt.Sprintf("type='signal',interface='fi.w1.wpa_supplicant1.Interface',path='%s'", self.Object.Path())
	call := self.WPA.Connection.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, match)
	self.Error = call.Err

	return self
}

func (self *InterfaceWPA) RemoveSignalsObserver() *InterfaceWPA {
	log.Log.Debug("RemoveSignalsObserver.Interface")

	match := fmt.Sprintf("type='signal',interface='fi.w1.wpa_supplicant1.Interface',path='%s'", self.Object.Path())
	call := self.WPA.Connection.BusObject().Call("org.freedesktop.DBus.RemoveMatch", 0, match)
	self.Error = call.Err

	return self
}

func (self *InterfaceWPA) ReadCurrentBSS() *InterfaceWPA {
	if self.Error != nil {
		return self
	}

	value, err := self.WPA.get("fi.w1.wpa_supplicant1.Interface.CurrentBSS", self.Object)
	if err != nil {
		self.Error = err
		return self
	}
	bssObjectPath := value.(dbus.ObjectPath)
	self.CurrentBSS = &BSSWPA{
		Interface: self,
		Object:    self.WPA.Connection.Object("fi.w1.wpa_supplicant1", bssObjectPath),
	}

	return self
}

func (self *InterfaceWPA) ReadCurrentNetwork() *InterfaceWPA {
	if self.Error != nil {
		return self
	}

	network, err := self.WPA.get("fi.w1.wpa_supplicant1.Interface.CurrentNetwork", self.Object)
	if err != nil {
		self.Error = err
		return self
	}
	networkObjectPath := network.(dbus.ObjectPath)
	self.CurrentNetwork = &NetworkWPA{Interface: self, Object: self.WPA.Connection.Object("fi.w1.wpa_supplicant1", networkObjectPath)}

	return self
}
