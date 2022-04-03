package wpa_dbus

import (
	"encoding/hex"
	"fmt"

	"github.com/godbus/dbus"
	"github.com/mark2b/wpa-connect/internal/log"
)

type BSSWPA struct {
	Interface     *InterfaceWPA
	Object        dbus.BusObject
	BSSID         string
	SSID          string
	WPAKeyMgmt    []string
	RSNKeyMgmt    []string
	WPS           string
	Frequency     uint16
	Signal        int16
	Age           uint32
	Mode          string
	Privacy       bool
	SignalChannel chan *dbus.Signal
	Error         error
}

func (self *BSSWPA) ReadWPA() *BSSWPA {
	if self.Error != nil {
		return self
	}

	variants, err := self.Interface.WPA.get("fi.w1.wpa_supplicant1.BSS.WPA", self.Object)
	if err != nil {
		self.Error = err
		return self
	}

	if variants, ok := variants.(map[string]dbus.Variant); ok {
		if keyMgmt, found := variants["KeyMgmt"]; found {
			self.RSNKeyMgmt = keyMgmt.Value().([]string)
		}
	}

	return self
}

func (self *BSSWPA) ReadRSN() *BSSWPA {
	if self.Error != nil {
		return self
	}

	variants, err := self.Interface.WPA.get("fi.w1.wpa_supplicant1.BSS.RSN", self.Object)
	if err != nil {
		self.Error = err
		return self
	}

	if variants, ok := variants.(map[string]dbus.Variant); ok {
		if keyMgmt, found := variants["KeyMgmt"]; found {
			self.RSNKeyMgmt = keyMgmt.Value().([]string)
		}
	}

	return self
}

func (self *BSSWPA) ReadWPS() *BSSWPA {
	if self.Error != nil {
		return self
	}

	variants, err := self.Interface.WPA.get("fi.w1.wpa_supplicant1.BSS.WPS", self.Object)
	if err != nil {
		self.Error = err
		return self
	}

	if variants, ok := variants.(map[string]dbus.Variant); ok {
		if wpsType, found := variants["Type"]; found {
			self.WPS = wpsType.String()
		}
	}

	return self
}

func (self *BSSWPA) ReadBSSID() *BSSWPA {
	if self.Error != nil {
		return self
	}

	value, err := self.Interface.WPA.get("fi.w1.wpa_supplicant1.BSS.BSSID", self.Object)
	if err != nil {
		self.Error = err
		return self
	}
	self.BSSID = hex.EncodeToString(value.([]byte))

	return self
}

func (self *BSSWPA) ReadSSID() *BSSWPA {
	if self.Error != nil {
		return self
	}

	value, err := self.Interface.WPA.get("fi.w1.wpa_supplicant1.BSS.SSID", self.Object)
	if err != nil {
		self.Error = err
		return self
	}
	self.SSID = string(value.([]byte))

	return self
}

func (self *BSSWPA) ReadFrequency() *BSSWPA {
	if self.Error != nil {
		return self
	}

	value, err := self.Interface.WPA.get("fi.w1.wpa_supplicant1.BSS.Frequency", self.Object)
	if err != nil {
		self.Error = err
		return self
	}
	self.Frequency = value.(uint16)

	return self
}

func (self *BSSWPA) ReadSignal() *BSSWPA {
	if self.Error != nil {
		return self
	}

	value, err := self.Interface.WPA.get("fi.w1.wpa_supplicant1.BSS.Signal", self.Object)
	if err != nil {
		self.Error = err
		return self
	}
	self.Signal = value.(int16)

	return self
}

func (self *BSSWPA) ReadAge() *BSSWPA {
	if self.Error != nil {
		return self
	}

	value, err := self.Interface.WPA.get("fi.w1.wpa_supplicant1.BSS.Age", self.Object)
	if err != nil {
		self.Error = err
		return self
	}
	self.Age = value.(uint32)

	return self
}

func (self *BSSWPA) ReadMode() *BSSWPA {
	if self.Error != nil {
		return self
	}

	value, err := self.Interface.WPA.get("fi.w1.wpa_supplicant1.BSS.Mode", self.Object)
	if err != nil {
		self.Error = err
		return self
	}
	self.Mode = value.(string)

	return self
}

func (self *BSSWPA) ReadPrivacy() *BSSWPA {
	if self.Error != nil {
		return self
	}

	value, err := self.Interface.WPA.get("fi.w1.wpa_supplicant1.BSS.Privacy", self.Object)
	if err != nil {
		self.Error = err
		return self
	}
	self.Privacy = value.(bool)

	return self
}

func (self *BSSWPA) AddSignalsObserver() *BSSWPA {
	log.Log.Debug("AddSignalsObserver.BSS")
	match := fmt.Sprintf("type='signal',interface='fi.w1.wpa_supplicant1.BSS',path='%s'", self.Object.Path())

	call := self.Interface.WPA.Connection.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, match)
	self.Error = call.Err

	return self
}

func (self *BSSWPA) RemoveSignalsObserver() *BSSWPA {
	log.Log.Debug("RemoveSignalsObserver.BSS")
	match := fmt.Sprintf("type='signal',interface='fi.w1.wpa_supplicant1.BSS',path='%s'", self.Object.Path())

	call := self.Interface.WPA.Connection.BusObject().Call("org.freedesktop.DBus.RemoveMatch", 0, match)
	self.Error = call.Err

	return self
}
