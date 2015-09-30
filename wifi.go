package main

import (
	"github.com/godbus/dbus"
)

func GetVisibleNetworks() ([]string, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		panic(err)
	}

	obj := conn.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager")

	networkDevices := []dbus.ObjectPath{}
	if err := obj.Call("org.freedesktop.NetworkManager.GetDevices", 0).Store(&networkDevices); err != nil {
		panic(err)
	}

	ret := []string{}

	for _, el := range networkDevices {
		obj := conn.Object("org.freedesktop.NetworkManager", el)
		accessPoints := []dbus.ObjectPath{}
		if err := obj.Call(
			"org.freedesktop.NetworkManager.Device.Wireless.GetAccessPoints",
			0,
		).Store(&accessPoints); err != nil {
			continue
		}
		for _, accessPoint := range accessPoints {
			obj := conn.Object("org.freedesktop.NetworkManager", accessPoint)

			ssid, err := obj.GetProperty(
				"org.freedesktop.NetworkManager.AccessPoint.Ssid",
			)

			if err != nil {
				panic(err)
			}
			ret = append(ret, string(ssid.Value().([]uint8)))
		}
	}

	return ret, nil
}
