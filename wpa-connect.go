package wpaconnect

import (
	"errors"

	"./internal/log"
	"./internal/wpa_cli"
	"./internal/wpa_dbus"
	"fmt"
	"github.com/godbus/dbus"
)

func (self *connectManager) Connect(ssid string, password string) (e error) {
	self.connectContext = &connectContext{}
	self.connectContext.scanDone = make(chan bool)
	self.connectContext.connectDone = make(chan bool)
	if wpa, err := wpa_dbus.NewWPA(); err == nil {
		wpa.WaifForSignals(self.onSignal)
		wpa.AddSignalsObserver()
		if wpa.ReadInterface(self.NetInterface); wpa.Error == nil {
			iface := wpa.Interface
			iface.AddSignalsObserver()
			self.connectContext.phaseWaitForScanDone = true
			if iface.Scan(); iface.Error == nil {
				// Wait for scan done
				<-self.connectContext.scanDone
				if iface.ReadBSSList(); iface.Error == nil {
					// Look for target BSS
					var bssFound = false
					for _, bss := range iface.BSSs {
						if bss.ReadSSID(); bss.Error == nil {
							log.Log.Info(bss.SSID)
							if ssid == bss.SSID {
								bssFound = true
								if err := self.connectToBSS(&bss, iface, password); err == nil {
									// Wait for connection
									cli := wpa_cli.WPACli{NetInterface: self.NetInterface}
									if err := cli.SaveConfig(); err == nil {
									} else {
										e = err
									}
								} else {
									e = err
								}
								break
							}
						} else {
							e = bss.Error
						}
					}
					if !bssFound {
						e = errors.New("ssid_not_found")
					}
				} else {
					e = iface.Error
				}
			} else {
				e = wpa.Error
			}
			iface.RemoveSignalsObserver()
		} else {
			e = wpa.Error
		}
		wpa.RemoveSignalsObserver()
		wpa.StopWaifForSignals()
	} else {
		e = err
	}
	log.Log.Debug("Connect exit")
	return
}

func (self *connectManager) connectToBSS(bss *wpa_dbus.BSSWPA, iface *wpa_dbus.InterfaceWPA, password string) (e error) {
	addNetworkArgs := map[string]dbus.Variant{
		"ssid": dbus.MakeVariant(bss.SSID),
		"psk":  dbus.MakeVariant(password)}
	if iface.RemoveAllNetworks().AddNetwork(addNetworkArgs); iface.Error == nil {
		network := iface.NewNetwork
		self.connectContext.phaseWaitForInterfaceConnected = true
		if network.Select(); network.Error == nil {
			connected := <-self.connectContext.connectDone
			log.Log.Debug("Connected", connected)
			if !connected {
				if iface.ReadDisconnectReason(); iface.Error == nil {
					e = errors.New(fmt.Sprintf("connection_failed, reason=%i", iface.DisconnectReason))
				} else {
					e = errors.New("connection_failed")
				}
			}
		} else {
			e = network.Error
		}
	} else {
		e = iface.Error
	}
	return
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

func (self *connectManager) processScanDone(wpa *wpa_dbus.WPA, signal *dbus.Signal) {
	log.Log.Debug("processScanDone")
	if self.connectContext.phaseWaitForScanDone {
		self.connectContext.phaseWaitForScanDone = false
		self.connectContext.scanDone <- true
	}
}

func (self *connectManager) processInterfacePropertiesChanged(wpa *wpa_dbus.WPA, signal *dbus.Signal) {
	log.Log.Debug("processInterfacePropertiesChanged")
	log.Log.Debug("phaseWaitForInterfaceConnected", self.connectContext.phaseWaitForInterfaceConnected)
	if self.connectContext.phaseWaitForInterfaceConnected {
		if len(signal.Body) > 0 {
			properties := signal.Body[0].(map[string]dbus.Variant)
			if stateVariant, hasState := properties["State"]; hasState {
				if state, ok := stateVariant.Value().(string); ok {
					log.Log.Debug("State", state)
					if state == "completed" {
						self.connectContext.phaseWaitForInterfaceConnected = false
						self.connectContext.connectDone <- true
						return
					} else if state == "disconnected" {
						self.connectContext.phaseWaitForInterfaceConnected = false
						self.connectContext.connectDone <- false
						return
					}
				}
			}
		}
	}
}

type connectContext struct {
	phaseWaitForScanDone           bool
	phaseWaitForInterfaceConnected bool
	scanDone                       chan bool
	connectDone                    chan bool
}

type connectManager struct {
	connectContext *connectContext
	NetInterface   string
}

var (
	ConnectManager = &connectManager{NetInterface: "wlan0"}
)
