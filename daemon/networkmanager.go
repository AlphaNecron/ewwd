package daemon

import (
	"github.com/nlepage/go-netmgr"
)

type (
	networkState struct {
		Icon        string `json:"icon"`
		Text        string `json:"text"`
		IsConnected bool   `json:"isConnected"`
	}
	accessPoint struct {
		ssid     string
		lastSeen int32
	}
)

var nwState = networkState{Icon: icons[1], IsConnected: true}

var icons = []string{"", "", "", ""}

func NetworkManager() {
	net, e := netmgr.System()
	if e != nil {
		panic(e)
	}
	/*getSsid := func() string {
		d, e := net.GetDeviceByIPIface("wlp2s0")
		if e != nil {
			return "Unknown"
		}
		r := d.Call("org.freedesktop.NetworkManager.Device.Wireless.GetAllAccessPoints", 0)
		if r.Err != nil || len(r.Body) == 0 {
			return "Unknown"
		}
		var aps []accessPoint
		for _, ap := range r.Body[0].([]dbus.ObjectPath) {
			obj := conn.Object(netmgr.BusName, ap)
			rawSsid, e := obj.GetProperty("org.freedesktop.NetworkManager.AccessPoint.Ssid")
			if e != nil {
				continue
			}
			lastSeen, e := obj.GetProperty("org.freedesktop.NetworkManager.AccessPoint.LastSeen")
			if e != nil {
				continue
			}
			aps = append(aps, accessPoint{string(rawSsid.Value().([]byte)), lastSeen.Value().(int32)})
		}
		if len(aps) == 0 {
			return "No connection"
		}
		sort.SliceStable(aps, func(i, j int) bool {
			return aps[i].lastSeen > aps[j].lastSeen
		})
		return aps[0].ssid
	}*/
	getSsid := func() string {
		c, e := net.PrimaryConnection()
		if e != nil {
			return "Disconnected"
		}
		v, e := c.GetProperty("org.freedesktop.NetworkManager.Connection.Active.Id")
		if e != nil {
			return "Unknown"
		}
		return v.Value().(string)
	}
	printNetworkState := func() {
		connState, _ := net.CheckConnectivity()
		nwState.IsConnected = true
		if connState == netmgr.ConnectivityNone {
			nwState.IsConnected = false
			nwState.Icon = icons[2]
			nwState.Text = "Disconnected"
		} else {
			nwState.Icon = icons[3]
			nwType, e := net.PrimaryConnectionType()
			if e == nil {
				if nwType == "802-11-wireless" {
					nwState.Icon = icons[1]
					nwState.Text = getSsid()
				} else {
					nwState.Icon = icons[0]
					nwState.Text = "Wired"
				}
			}
			if connState != netmgr.ConnectivityFull {
				nwState.Icon = icons[3]
			}
		}
		stdout.Encode(nwState)
	}
	state := make(chan netmgr.StateEnum)
	if e := net.StateChanged(state); e != nil {
		panic(e)
	}
	printNetworkState()
	for range state {
		printNetworkState()
	}
}
