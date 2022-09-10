package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/loafoe/prometheus-p1-exporter/conn"
	"github.com/loafoe/prometheus-p1-exporter/parser"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
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

func main() {
	flag.StringVar(&listenAddr, "listen", "127.0.0.1:8888", "Listen address for HTTP metrics")
	flag.DurationVar(&readInterval, "interval", 10*time.Second, "Interval between metric reads")
	flag.BoolVar(&useMock, "mock", false, "Use dummy source instead of ttyUSB0 socket")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output logging")
	flag.Parse()

	var source conn.Source
	if useMock {
		source = &conn.MockSource{}
	} else {
		source = &conn.SerialSource{}
	}

	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	go func() {
		errorCount := 0
		for {
			if errorCount > 10 {
				logrus.Errorln("Quitting because there were too many errors")
				os.Exit(1)
			}

			lines, err := conn.ReadTelegram(&conn.ESMR5TelegramReaderOptions, source)
			if err != nil {
				logrus.Errorln("Error while reading telegram from source", err)
				errorCount++
				time.Sleep(readInterval)
				continue
			}
			telegram, err := parser.ParseTelegram(&parser.XS210ESMR5TelegramFormat, lines)
			if err != nil {
				logrus.Errorln("Error while parsing telegram", err)
				errorCount++
				time.Sleep(readInterval)
				continue
			}
			errorCount = 0
			electricityUsageHighMetric.Set(telegram.ElectricityUsageHigh)
			electricityUsageLowMetric.Set(telegram.ElectricityUsageLow)
			electricityReturnedHighMetric.Set(telegram.ElectricityReturnedHigh)
			electricityReturnedLowMetric.Set(telegram.ElectricityReturnedLow)
			actualElectricityDeliveredMetric.Set(telegram.ActiveElectricityDraw)
			actualElectricityRetreivedMetric.Set(telegram.ActiveElectricityReturn)
			activeTarrifMetric.Set(float64(telegram.ActiveTariff))
			powerFailuresLongMetric.Set(float64(telegram.PowerFailuresLong))
			powerFailuresShortMetric.Set(float64(telegram.PowerFailuresShort))
			gasUsageMetric.Set(telegram.GasUsage)

			logrus.Debugf("%+v\n", telegram)

			time.Sleep(readInterval)
		}
	}()

	logrus.Infoln("Start listening at", listenAddr)
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	logrus.Fatalln(http.ListenAndServe(listenAddr, nil))
}
