# wpa-connect
## Package provides API for connection Linux device to Wi-Fi Network.


**wpa-connect** communicates with WPA supplicant over D-Bus (linux message bus system).


This package was developed as part of IoT project in order to add Wi-Fi connectivity to headless Raspberry Pi like devices. No need for **connman** or **Network Manager** be installed. 


## Setup

**On Linux:**

**wpa_supplicant** service should run with **-u** flag in order to enable DBus interface. Run it as Linux service before first call to **wpa_supplicant**. Otherwise system will start it automatically without **-u** flag. 

Systemd service configuration file - **/etc/systemd/system/wpa_supplicant@wlan0.service**
```
[Unit]
Description=WPA supplicant for %i

[Service]
ExecStart=/usr/sbin/wpa_supplicant -u -i%i -c/etc/wpa_supplicant.conf -Dwext

[Install]
WantedBy=multi-user.target
```

**On Project:**

```
go get github.com/mark2b/wpa-connect
```

## Usage
Please see [godoc.org](http://godoc.org/github.com/mark2b/wpa-connect) for documentation. (Not ready yet)

## Examples

### Connect to Wi-Fi network

 
```golang
import wifi "github.com/mark2b/wpa-connect"

if conn, err := wifi.ConnectManager.Connect(ssid, password, time.Second * 60); err == nil {
	fmt.Println("Connected", conn.NetInterface, conn.SSID, conn.IP4.String(), conn.IP6.String())
} else {
	fmt.Println(err)
}
```
### Scan for Wi-Fi networks

```golang
import wifi "github.com/mark2b/wpa-connect"

if bssList, err := wifi.ScanManager.Scan(); err == nil {
	for _, bss := range bssList {
		print(bss.SSID, bss.Signal, bss.KeyMgmt)
	}
}
```

Package release under a [MIT license](./LICENSE.md).
