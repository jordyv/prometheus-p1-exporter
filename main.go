package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/jordyv/prometheus-p1-exporter/conn"
	"github.com/jordyv/prometheus-p1-exporter/parser"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var readInterval time.Duration
var listenAddr string
var apiEndpoint string
var useMock bool
var verbose bool
var metricNamePrefix = "p1_"
var usbSerial string

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
	flag.StringVar(&apiEndpoint, "apiEndpoint", "", "Use API endpoint to read the telegram (use for HomeWizard)")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output logging")
	flag.StringVar(&usbSerial, "usbserial", "/dev/ttyUSB0", "USB serial device")
	flag.Parse()

	var source conn.Source
	if useMock {
		source = conn.NewMockSource()
	} else if apiEndpoint != "" {
		source = conn.NewAPISource(apiEndpoint)
	} else {
		source = conn.NewSerialSource(usbSerial)
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
			if telegram.ElectricityUsageHigh != nil {
				electricityUsageHighMetric.Set(*telegram.ElectricityUsageHigh)
			}
			if telegram.ElectricityUsageLow != nil {
				electricityUsageLowMetric.Set(*telegram.ElectricityUsageLow)
			}
			if telegram.ElectricityReturnedHigh != nil {
				electricityReturnedHighMetric.Set(*telegram.ElectricityReturnedHigh)
			}
			if telegram.ElectricityReturnedLow != nil {
				electricityReturnedLowMetric.Set(*telegram.ElectricityReturnedLow)
			}
			if telegram.ActualElectricityDelivered != nil {
				actualElectricityDeliveredMetric.Set(*telegram.ActualElectricityDelivered)
			}
			if telegram.ActualElectricityRetreived != nil {
				actualElectricityRetreivedMetric.Set(*telegram.ActualElectricityRetreived)
			}
			activeTarrifMetric.Set(float64(telegram.ActiveTariff))
			powerFailuresLongMetric.Set(float64(telegram.PowerFailuresLong))
			powerFailuresShortMetric.Set(float64(telegram.PowerFailuresShort))
			if telegram.GasUsage != nil {
				gasUsageMetric.Set(*telegram.GasUsage)
			}

			logrus.Debugf("%+v\n", telegram)

			time.Sleep(readInterval)
		}
	}()

	logrus.Infoln("Start listening at", listenAddr)
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	logrus.Fatalln(http.ListenAndServe(listenAddr, nil))
}
