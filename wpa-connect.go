package wpaconnect

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/mark2b/wpa-connect/internal/wpa_cli"

	"github.com/godbus/dbus"
	"github.com/mark2b/wpa-connect/internal/log"
	"github.com/mark2b/wpa-connect/internal/wpa_dbus"
)

type ConnectionInfo struct {
	NetInterface string
	SSID         string
	IP4          net.IP
	IP6          net.IP
}

type connectContext struct {
	phaseWaitForScanDone           bool
	phaseWaitForInterfaceConnected bool
	scanDone                       chan bool
	connectDone                    chan bool
	ip4                            net.IP
	ip6                            net.IP
	err                            error
}

type connectManager struct {
	context      *connectContext
	deadTime     time.Time
	NetInterface string
}

var ConnectManager = &connectManager{NetInterface: "wlan0"}

func NewConnectManager(netInterface string) *connectManager {
	return &connectManager{NetInterface: netInterface}
}

func (self *connectManager) Connect(ssid string, password string, timeout time.Duration) (ConnectionInfo, error) {
	self.deadTime = time.Now().Add(timeout)
	self.context = &connectContext{}
	self.context.scanDone = make(chan bool)
	self.context.connectDone = make(chan bool)

	wpa, err := wpa_dbus.NewWPA()
	if err != nil {
		return ConnectionInfo{}, err
	}

	wpa.WaitForSignals(self.onSignal)
	defer wpa.StopWaitForSignals()

	wpa.AddSignalsObserver()
	defer wpa.RemoveSignalsObserver()

	wpa.ReadInterface(self.NetInterface)
	if wpa.Error != nil {
		return ConnectionInfo{}, wpa.Error
	}

	iface := wpa.Interface
	iface.AddSignalsObserver()
	defer iface.RemoveSignalsObserver()

	self.context.phaseWaitForScanDone = true
	go func() {
		time.Sleep(time.Until(self.deadTime))
		self.context.scanDone <- false
		self.context.err = errors.New("timeout")
	}()

	iface.Scan()
	if iface.Error != nil {
		return ConnectionInfo{}, iface.Error
	}

	// Wait for scan done
	<-self.context.scanDone
	if self.context.err != nil {
		return ConnectionInfo{}, self.context.err
	}

	iface.ReadBSSList()
	if iface.Error != nil {
		return ConnectionInfo{}, iface.Error
	}

	bssMap := make(map[string]wpa_dbus.BSSWPA, 0)
	for _, bss := range iface.BSSs {
		bss.ReadSSID()
		if bss.Error != nil {
			return ConnectionInfo{}, bss.Error
		}

		bssMap[bss.SSID] = bss
		log.Log.Debug(bss.SSID, bss.BSSID)
	}

	_, exists := bssMap[ssid]
	err = self.connectToBSS(&wpa_dbus.BSSWPA{
		SSID: ssid,
	}, iface, password, !exists)
	if err != nil {
		return ConnectionInfo{}, err
	}

	// Connected, save configuration
	cli := wpa_cli.WPACli{NetInterface: self.NetInterface}
	err = cli.SaveConfig()
	if err != nil {
		return ConnectionInfo{}, err
	}

	return ConnectionInfo{
		NetInterface: self.NetInterface,
		SSID:         ssid,
		IP4:          self.context.ip4,
		IP6:          self.context.ip6,
	}, nil
}

func (self *connectManager) connectToBSS(bss *wpa_dbus.BSSWPA, iface *wpa_dbus.InterfaceWPA, password string, isHidden bool) error {
	addNetworkArgs := map[string]dbus.Variant{
		"ssid": dbus.MakeVariant(bss.SSID),
	}

	if isHidden {
		addNetworkArgs["scan_ssid"] = dbus.MakeVariant(1)
	}

	if password == "" {
		addNetworkArgs["key_mgmt"] = dbus.MakeVariant("NONE")
	} else {
		addNetworkArgs["psk"] = dbus.MakeVariant(password)
	}

	iface.RemoveAllNetworks().AddNetwork(addNetworkArgs)
	if iface.Error != nil {
		return iface.Error
	}

	network := iface.NewNetwork
	self.context.phaseWaitForInterfaceConnected = true
	go func() {
		time.Sleep(time.Until(self.deadTime))
		self.context.connectDone <- false
		self.context.err = errors.New("timeout")
	}()

	network.Select()
	if network.Error != nil {
		return network.Error
	}

	connected := <-self.context.connectDone
	if self.context.err != nil {
		return self.context.err
	}

	if connected {
		err := self.readNetAddress()
		if err != nil {
			return err
		}
	} else {
		iface.ReadDisconnectReason()
		if iface.Error != nil {
			return fmt.Errorf("connection failed, reason=%d", iface.DisconnectReason)
		}
		return errors.New("connection failed")
	}

	return nil
}

func (self *connectManager) onSignal(wpa *wpa_dbus.WPA, signal *dbus.Signal) {
	log.Log.Debug(signal.Name, signal.Path)
	switch signal.Name {
	case "fi.w1.wpa_supplicant1.Interface.BSSAdded":
	case "fi.w1.wpa_supplicant1.Interface.BSSRemoved":
		break
	case "fi.w1.wpa_supplicant1.Interface.ScanDone":
		self.processScanDone(wpa, signal)
	case "fi.w1.wpa_supplicant1.Interface.PropertiesChanged":
		log.Log.Debug(signal.Name, signal.Path, signal.Body)
		self.processInterfacePropertiesChanged(wpa, signal)
	default:
		log.Log.Debug(signal.Name, signal.Path, signal.Body)
	}
}

func (self *connectManager) readNetAddress() error {
	netIface, err := net.InterfaceByName(self.NetInterface)
	if err != nil {
		return err
	}

	for time.Now().Before(self.deadTime) && !self.context.hasIP() {
		addrs, _ := netIface.Addrs()
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if self.context.ip4 == nil {
					self.context.ip4 = ipnet.IP.To4()
					continue
				}
				if self.context.ip6 == nil {
					self.context.ip6 = ipnet.IP.To16()
				}
			}
		}

		time.Sleep(time.Millisecond * 500)
	}
	if !self.context.hasIP() {
		return errors.New("address not allocated")
	}

	return nil
}

func (self *connectManager) processScanDone(wpa *wpa_dbus.WPA, signal *dbus.Signal) {
	log.Log.Debug("processScanDone")
	if self.context.phaseWaitForScanDone {
		self.context.phaseWaitForScanDone = false
		self.context.scanDone <- true
	}
}

func (self *connectManager) processInterfacePropertiesChanged(wpa *wpa_dbus.WPA, signal *dbus.Signal) {
	log.Log.Debug("processInterfacePropertiesChanged")
	log.Log.Debug("phaseWaitForInterfaceConnected", self.context.phaseWaitForInterfaceConnected)
	if !self.context.phaseWaitForInterfaceConnected {
		return
	}

	if len(signal.Body) == 0 {
		return
	}

	properties := signal.Body[0].(map[string]dbus.Variant)
	if stateVariant, hasState := properties["State"]; hasState {
		if state, ok := stateVariant.Value().(string); ok {
			log.Log.Debug("State", state)
			switch state {
			case "completed":
				self.context.phaseWaitForInterfaceConnected = false
				self.context.connectDone <- true
			case "disconnected":
				// self.context.phaseWaitForInterfaceConnected = false
				// self.context.connectDone <- false
			}
		}
	}
}

func (self *connectContext) hasIP() bool {
	return self.ip4 != nil && self.ip6 != nil
}
