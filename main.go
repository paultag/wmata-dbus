package main

import (
	"fmt"
	"log"
	"os"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"

	"pault.ag/go/config"
	"pault.ag/go/wmata"
)

type WMATADbusInterface struct{}

func (w WMATADbusInterface) NextTrains(stops []string) ([]map[string]string, *dbus.Error) {
	log.Printf("Getting info")
	predictions, err := wmata.GetPredictionsByCodes(stops...)
	if err != nil {
		return []map[string]string{}, dbus.NewError(
			"org.anized.wmata.Rail.NotFound",
			[]interface{}{err.Error()},
		)
	}

	log.Printf("Building map")
	ret := []map[string]string{}
	for _, prediction := range predictions {
		ret = append(ret, map[string]string{
			"cars":             prediction.Cars,
			"group":            prediction.Group,
			"line":             prediction.Line.Code,
			"minutes":          prediction.Minutes,
			"desitnation":      prediction.Destination,
			"desitnation_name": prediction.DesitnationName,
			"desitnation_code": prediction.DesitnationCode,
			"location_name":    prediction.LocationName,
			"location_code":    prediction.LocationCode,
		})
	}
	return ret, nil
}

type WMATADbus struct {
	APIKey string `flag:"apikey" description:"API Key to use"`
}

func main() {
	conf := WMATADbus{}
	if err := config.Load("wmatadbusd", &conf); err != nil {
		panic(err)
	}
	wmata.SetAPIKey(conf.APIKey)

	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}
	reply, err := conn.RequestName("org.anized.wmata.Rail",
		dbus.NameFlagDoNotQueue)
	if err != nil {
		panic(err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Fprintln(os.Stderr, "name already taken")
		os.Exit(1)
	}

	wmata := WMATADbusInterface{}
	introspectedMethods := introspect.Methods(wmata)

	node := introspect.Node{
		Name: "/org/anized/wmata",
		Interfaces: []introspect.Interface{
			introspect.Interface{
				Name:    "org.anized.wmata.Rail",
				Methods: introspectedMethods,
			},
		},
	}

	export := introspect.NewIntrospectable(&node)
	conn.Export(wmata, "/org/anized/wmata/Rail", "org.anized.wmata.Rail")
	conn.Export(
		export,
		"/org/anized/wmata/Rail",
		"org.freedesktop.DBus.Introspectable",
	)
	select {}
}
