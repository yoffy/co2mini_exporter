package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zserge/hid"
)

const (
	vendorID = "04d9:a052:0200:00"
	co2op    = 0x50
	tempop   = 0x42
)

var (
	key = []byte{0x86, 0x41, 0xc9, 0xa8, 0x7f, 0x41, 0x3c, 0xac}
)

var (
	co2Gauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "co2mini_co2",
		Help: "co2 in ppm",
	})
	temperatureGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "co2mini_temperature",
		Help: "temperature in Celsius",
	})
)

// protocol: http://co2meters.com/Documentation/AppNotes/AN146-RAD-0401-serial-communication.pdf
func monitor(device hid.Device) {
	if err := device.Open(); err != nil {
		log.Println("Open error: ", err)
		return
	}
	defer device.Close()

	if err := device.SetReport(0, key); err != nil {
		log.Fatal(err)
	}

	for {
		if buf, err := device.Read(-1, 1*time.Second); err == nil {
			if len(buf) == 0 {
				continue
			}
			val := int(buf[1])<<8 | int(buf[2])
			if buf[0] == co2op {
				co2Gauge.Set(float64(val))
			}
			if buf[0] == tempop {
				temp := float64(val)/16.0 - 273.15
				temperatureGauge.Set(temp)
			}
		}
	}
}

func main() {
	// define metrics
	prometheus.MustRegister(co2Gauge)
	prometheus.MustRegister(temperatureGauge)
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())

	// find CO2-MINI
	hid.UsbWalk(func(device hid.Device) {
		info := device.Info()
		id := fmt.Sprintf("%04x:%04x:%04x:%02x", info.Vendor, info.Product, info.Revision, info.Interface)
		if id != vendorID {
			return
		}

		// start minitoring CO2-MINI
		go monitor(device)
	})

	// start HTTP server
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9002", nil))
}
