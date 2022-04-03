package wpa_dbus

import (
	"errors"
	"fmt"

	"github.com/godbus/dbus"
	"github.com/mark2b/wpa-connect/internal/log"
)

type WPA struct {
	Connection    *dbus.Conn
	Object        dbus.BusObject
	Interfaces    []InterfaceWPA
	Interface     *InterfaceWPA
	SignalChannel chan *dbus.Signal
	Error         error
}

func NewWPA() (*WPA, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	if obj := conn.Object("fi.w1.wpa_supplicant1", "/fi/w1/wpa_supplicant1"); obj != nil {
		return &WPA{Connection: conn, Object: obj}, nil
	}

	conn.Close()
	return nil, errors.New("Can't create WPA object")
}

func (self *WPA) ReadInterface(ifname string) *WPA {
	if self.Error != nil {
		return self
	}

	call := self.Object.Call("fi.w1.wpa_supplicant1.GetInterface", 0, ifname)
	if call.Err != nil {
		self.Error = call.Err
		return self
	}

	objectPath := dbus.ObjectPath(call.Body[0].(dbus.ObjectPath))
	self.Interface = &InterfaceWPA{
		WPA:    self,
		Object: self.Connection.Object("fi.w1.wpa_supplicant1", objectPath),
	}

	return self
}

func (self *WPA) ReadInterfaceList() *WPA {
	if self.Error != nil {
		return self
	}

	interfaces, err := self.get("fi.w1.wpa_supplicant1.Interfaces", self.Object)
	if err != nil {
		self.Error = err
		return self
	}

	newInterfaces := []InterfaceWPA{}
	for _, interfaceObjectPath := range interfaces.([]dbus.ObjectPath) {
		iface := InterfaceWPA{
			WPA:    self,
			Object: self.Connection.Object("fi.w1.wpa_supplicant1", interfaceObjectPath),
		}
		newInterfaces = append(newInterfaces, iface)
	}
	self.Interfaces = newInterfaces

	return self
}

func (self *WPA) get(name string, target dbus.BusObject) (interface{}, error) {
	variant, err := target.GetProperty(name)
	if err != nil {
		return nil, err
	}

	return variant.Value(), nil
}

func (self *WPA) WaitForSignals(callBack func(*WPA, *dbus.Signal)) *WPA {
	log.Log.Debug("WaitForSignals")

	self.SignalChannel = make(chan *dbus.Signal, 10)
	self.Connection.Signal(self.SignalChannel)
	go func() {
		for ch := range self.SignalChannel {
			callBack(self, ch)
		}
	}()

	return self
}

func (self *WPA) StopWaitForSignals() *WPA {
	log.Log.Debug("StopWaitForSignals")

	self.Connection.RemoveSignal(self.SignalChannel)

	return self
}

func (self *WPA) AddSignalsObserver() *WPA {
	log.Log.Debug("AddSignalsObserver.WPA")

	match := fmt.Sprintf("type='signal',interface='fi.w1.wpa_supplicant1',path='%s'", self.Object.Path())
	call := self.Connection.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, match)
	self.Error = call.Err

	return self
}

func (self *WPA) RemoveSignalsObserver() *WPA {
	log.Log.Debug("RemoveSignalsObserver.WPA")

	match := fmt.Sprintf("type='signal',interface='fi.w1.wpa_supplicant1',path='%s'", self.Object.Path())
	call := self.Connection.BusObject().Call("org.freedesktop.DBus.RemoveMatch", 0, match)
	self.Error = call.Err

	return self
}
