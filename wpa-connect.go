package wpaconnect

import (
	"errors"

	"fmt"
	"github.com/godbus/dbus"
	"github.com/mark2b/wpa-connect/internal/log"
	"github.com/mark2b/wpa-connect/internal/wpa_cli"
	"github.com/mark2b/wpa-connect/internal/wpa_dbus"
	"net"
	"time"
)

func (self *connectManager) Connect(ssid string, password string, timeout time.Duration) (connectionInfo ConnectionInfo, e error) {
	self.deadTime = time.Now().Add(timeout)
	self.context = &connectContext{}
	self.context.scanDone = make(chan bool)
	self.context.connectDone = make(chan bool)
	if wpa, err := wpa_dbus.NewWPA(); err == nil {
		wpa.WaitForSignals(self.onSignal)
		wpa.AddSignalsObserver()
		if wpa.ReadInterface(self.NetInterface); wpa.Error == nil {
			iface := wpa.Interface
			iface.AddSignalsObserver()
			self.context.phaseWaitForScanDone = true
			go func() {
				time.Sleep(self.deadTime.Sub(time.Now()))
				self.context.scanDone <- false
				self.context.error = errors.New("timeout")
			}()
			if iface.Scan(); iface.Error == nil {
				// Wait for scan_example done
				if <-self.context.scanDone; self.context.error == nil {
					if iface.ReadBSSList(); iface.Error == nil {
						// Look for target BSS
						var bssFound = false
						for _, bss := range iface.BSSs {
							if bss.ReadSSID(); bss.Error == nil {
								if ssid == bss.SSID {
									bssFound = true
									if err := self.connectToBSS(&bss, iface, password); err == nil {
										// Wait for connection
										cli := wpa_cli.WPACli{NetInterface: self.NetInterface}
										if err := cli.SaveConfig(); err == nil {
											connectionInfo = ConnectionInfo{NetInterface: self.NetInterface, SSID: ssid,
												IP4: self.context.ip4, IP6: self.context.ip6}
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
					e = self.context.error
				}
			} else {
				e = wpa.Error
			}
			iface.RemoveSignalsObserver()
		} else {
			e = wpa.Error
		}
		wpa.RemoveSignalsObserver()
		wpa.StopWaitForSignals()
	} else {
		e = err
	}
	return
}

func (self *connectManager) Disconnect() (e error) {
	if wpa, err := wpa_dbus.NewWPA(); err == nil {
		if wpa.ReadInterface(self.NetInterface); wpa.Error == nil {
			iface := wpa.Interface
			if iface.RemoveAllNetworks(); iface.Error == nil {
				e = nil
			}
		} else {
			e = wpa.Error
		}
	} else {
		e = err
	}
	return
}

func (self *connectManager) connectToBSS(bss *wpa_dbus.BSSWPA, iface *wpa_dbus.InterfaceWPA, password string) (e error) {
	addNetworkArgs := map[string]dbus.Variant{
		"ssid": dbus.MakeVariant(bss.SSID),
		"psk":  dbus.MakeVariant(password)}
	if iface.AddNetwork(addNetworkArgs); iface.Error == nil {
		network := iface.NewNetwork
		self.context.phaseWaitForInterfaceConnected = true
		go func() {
			time.Sleep(self.deadTime.Sub(time.Now()))
			self.context.connectDone <- false
			self.context.error = errors.New("timeout")
		}()
		if network.Select(); network.Error == nil {
			if connected := <-self.context.connectDone; self.context.error == nil {
				if connected {
					if err := self.readNetAddress(); err == nil {
					} else {
						e = err
					}
				} else {
					if iface.ReadDisconnectReason(); iface.Error == nil {
						e = errors.New(fmt.Sprintf("connection_failed, reason=%d", iface.DisconnectReason))
					} else {
						e = errors.New("connection_failed")
					}
				}
			} else {
				e = self.context.error
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

func (self *connectManager) readNetAddress() (e error) {
	if netIface, err := net.InterfaceByName(self.NetInterface); err == nil {
		for time.Now().Before(self.deadTime) && !self.context.hasIP() {
			if addrs, err := netIface.Addrs(); err == nil {
				for _, addr := range addrs {
					if ip, _, err := net.ParseCIDR(addr.String()); err == nil {
						if self.context.ip4 == nil {
							self.context.ip4 = ip.To4()
							continue
						}
						if self.context.ip6 == nil {
							self.context.ip6 = ip.To16()
							continue
						}
					} else {
						e = err
						return
					}
				}
			} else {
				e = err
			}
			time.Sleep(time.Millisecond * 500)
		}
		if !self.context.hasIP() {
			e = errors.New("address_not_allocated")
		}
	} else {
		e = err
	}
	return
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
	if self.context.phaseWaitForInterfaceConnected {
		if len(signal.Body) > 0 {
			properties := signal.Body[0].(map[string]dbus.Variant)
			if stateVariant, hasState := properties["State"]; hasState {
				if state, ok := stateVariant.Value().(string); ok {
					log.Log.Debug("State", state)
					if state == "completed" {
						self.context.phaseWaitForInterfaceConnected = false
						self.context.connectDone <- true
						return
					} else if state == "disconnected" {
						self.context.phaseWaitForInterfaceConnected = false
						self.context.connectDone <- false
						return
					}
				}
			}
		}
	}
}

func (self *connectContext) hasIP() bool {
	return self.ip4 != nil && self.ip6 != nil
}

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
	error                          error
}

type connectManager struct {
	context      *connectContext
	deadTime     time.Time
	NetInterface string
}

var (
	ConnectManager = &connectManager{NetInterface: "wlan0"}
)
