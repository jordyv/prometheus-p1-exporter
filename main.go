package main

import (
	"flag"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/skoef/gop1"
)

var readInterval time.Duration
var listenAddr string
var useMock bool
var verbose bool
var metricNamePrefix = "p1_"

var (
	registry                   = prometheus.NewRegistry()
	electricityUsageHighMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metricNamePrefix + "usage_electricity_high",
		Help: "Electricity usage high tariff",
	})
	electricityUsageLowMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metricNamePrefix + "usage_electricity_low",
		Help: "Electricity usage low tariff",
	})
	electricityReturnedHighMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metricNamePrefix + "returned_electricity_high",
		Help: "Electricity returned high tariff",
	})
	electricityReturnedLowMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metricNamePrefix + "returned_electricity_low",
		Help: "Electricity returned low tariff",
	})
	actualElectricityDeliveredMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metricNamePrefix + "actual_electricity_delivered",
		Help: "Actual electricity power delivered to client",
	})
	actualElectricityRetreivedMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metricNamePrefix + "actual_electricity_retreived",
		Help: "Actual electricity power retreived from client",
	})
	activeTarrifMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metricNamePrefix + "active_tariff",
		Help: "Active tariff",
	})
	powerFailuresLongMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metricNamePrefix + "power_failures_long",
		Help: "Power failures long",
	})
	powerFailuresShortMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metricNamePrefix + "power_failures_short",
		Help: "Power failures short",
	})
	gasUsageMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metricNamePrefix + "usage_gas",
		Help: "Gas usage",
	})
)

func init() {
	registry.MustRegister(electricityUsageHighMetric)
	registry.MustRegister(electricityUsageLowMetric)
	registry.MustRegister(electricityReturnedHighMetric)
	registry.MustRegister(electricityReturnedLowMetric)
	registry.MustRegister(actualElectricityDeliveredMetric)
	registry.MustRegister(actualElectricityRetreivedMetric)
	registry.MustRegister(activeTarrifMetric)
	registry.MustRegister(powerFailuresLongMetric)
	registry.MustRegister(powerFailuresShortMetric)
	registry.MustRegister(gasUsageMetric)
}

func floatValue(input string) (fval float64) {
	fval, _ = strconv.ParseFloat(input, 64)
	return
}

func main() {
	flag.StringVar(&listenAddr, "listen", "127.0.0.1:8888", "Listen address for HTTP metrics")
	flag.DurationVar(&readInterval, "interval", 10*time.Second, "Interval between metric reads")
	flag.BoolVar(&useMock, "mock", false, "Use dummy source instead of ttyUSB0 socket")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output logging")
	flag.Parse()

	p1, err := gop1.New(gop1.P1Config{
		USBDevice: "/dev/ttyUSB0",
	})
	if err != nil {
		logrus.Errorln("Quitting because of error opening p1", err)
		os.Exit(1)
	}

	// Start
	p1.Start()

	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	go func() {
		for tgram := range p1.Incoming {
			for _, obj := range tgram.Objects {
				switch obj.Type {

				case gop1.OBISTypeElectricityDelivered:
					actualElectricityDeliveredMetric.Set(floatValue(obj.Values[0].Value))
				case gop1.OBISTypeElectricityGenerated:
					actualElectricityRetreivedMetric.Set(floatValue(obj.Values[0].Value))
				}
			}
		}
	}()

	logrus.Infoln("Start listening at", listenAddr)
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	logrus.Fatalln(http.ListenAndServe(listenAddr, nil))
}
