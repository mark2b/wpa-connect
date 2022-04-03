package wpaconnect

import (
	"github.com/godbus/dbus"
	"github.com/mark2b/wpa-connect/internal/log"
	"github.com/mark2b/wpa-connect/internal/wpa_dbus"
)

type BSS struct {
	BSSID     string
	SSID      string
	KeyMgmt   []string
	WPS       string
	Frequency uint16
	Signal    int16
	Age       uint32
	Mode      string
	Privacy   bool
}

type scanContext struct {
	phaseWaitForScanDone bool
	scanDone             chan bool
}

type scanManager struct {
	scanContext  *scanContext
	NetInterface string
}

var ScanManager = &scanManager{NetInterface: "wlan0"}

func NewScanManager(netInterface string) *scanManager {
	return &scanManager{NetInterface: netInterface}
}

func (self *scanManager) Scan() ([]BSS, error) {
	self.scanContext = &scanContext{}
	self.scanContext.scanDone = make(chan bool)

	wpa, err := wpa_dbus.NewWPA()
	if err != nil {
		return nil, err
	}

	wpa.WaitForSignals(self.onScanSignal)
	defer wpa.StopWaitForSignals()

	wpa.ReadInterface(self.NetInterface)
	if wpa.Error != nil {
		return nil, wpa.Error
	}

	iface := wpa.Interface
	iface.AddSignalsObserver()
	defer iface.RemoveSignalsObserver()
	self.scanContext.phaseWaitForScanDone = true

	iface.Scan()
	if iface.Error != nil {
		return nil, iface.Error
	}

	// Wait for scan_example done
	<-self.scanContext.scanDone
	bssList := []BSS{}
	iface.ReadBSSList()
	if iface.Error != nil {
		return nil, iface.Error
	}

	for _, bss := range iface.BSSs {
		bss.ReadBSSID().ReadSSID().ReadRSN().ReadMode().ReadSignal().
			ReadFrequency().ReadPrivacy().ReadAge().ReadWPS().ReadWPA()
		if bss.Error != nil {
			continue
		}

		bssList = append(bssList, BSS{
			BSSID: bss.BSSID, SSID: bss.SSID, KeyMgmt: bss.RSNKeyMgmt, WPS: bss.WPS,
			Frequency: bss.Frequency, Privacy: bss.Privacy, Age: bss.Age, Mode: bss.Mode, Signal: bss.Signal,
		})
	}

	return bssList, nil
}

func (self *scanManager) onScanSignal(wpa *wpa_dbus.WPA, signal *dbus.Signal) {
	log.Log.Debug(signal.Name, signal.Path)
	switch signal.Name {
	case "fi.w1.wpa_supplicant1.Interface.BSSAdded":
	case "fi.w1.wpa_supplicant1.Interface.BSSRemoved":
	case "fi.w1.wpa_supplicant1.Interface.PropertiesChanged":
		break
	case "fi.w1.wpa_supplicant1.Interface.ScanDone":
		self.processScanDone(wpa, signal)
	default:
		log.Log.Debug(signal.Name, signal.Path, signal.Body)
	}
}

func (self *scanManager) processScanDone(wpa *wpa_dbus.WPA, signal *dbus.Signal) {
	log.Log.Debug("processScanDone")
	if self.scanContext.phaseWaitForScanDone {
		self.scanContext.phaseWaitForScanDone = false
		self.scanContext.scanDone <- true
	}
}
