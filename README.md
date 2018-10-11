# Prometheus P1 exporter #

Prometheus exporter for smart meter statistics fetched with a P1 cable.

## Todos ##

[] Add build target for a .deb package to the Makefile
[] Add more unit tests
[] Fix TravisCI integration

## Installation ##

### Debian/Ubuntu ###

TODO

### From source ###

With Go get:

```
$ go get github.com/jordyv/prometheus-p1-exporter
```

Make:

```
$ git clone https://github.com/jordyv/prometheus-p1-exporter.git
$ cd prometheus-p1-exporter
$ make
```

## Usage ##

```
Usage of ./prometheus-p1-exporter:
  -interval duration
        Interval between metric reads (default 10s)
  -listen string
        Listen address for HTTP metrics (default "127.0.0.1:8888")
  -mock
        Use dummy source instead of ttyUSB0 socket
  -verbose
        Verbose output logging
```

By default the exporter will collect metrics from `/dev/ttyUSB0` every 10 seconds and export the metrics to an HTTP endpoint at `http://127.0.0.1:8888/metrics`. This endpoint can be added to your Prometheus configuration.

Example metrics page:

```
# HELP active_tariff Active tariff
# TYPE active_tariff gauge
active_tariff 2
# HELP current_usage_electricity_high Electricity currently used high tariff
# TYPE current_usage_electricity_high gauge
current_usage_electricity_high 0
# HELP current_usage_electricity_low Electricity currently used low tariff
# TYPE current_usage_electricity_low gauge
current_usage_electricity_low 0.2
# HELP power_failures_long Power failures long
# TYPE power_failures_long gauge
power_failures_long 2
# HELP power_failures_short Power failures short
# TYPE power_failures_short gauge
power_failures_short 57
# HELP returned_electricity_high Electricity returned high tariff
# TYPE returned_electricity_high gauge
returned_electricity_high 0
# HELP returned_electricity_low Electricity returned low tariff
# TYPE returned_electricity_low gauge
returned_electricity_low 0.016
# HELP usage_electricity_high Electricity usage high tariff
# TYPE usage_electricity_high gauge
usage_electricity_high 1225.59
# HELP usage_electricity_low Electricity usage low tariff
# TYPE usage_electricity_low gauge
usage_electricity_low 1179.186
# HELP usage_gas Gas usage
# TYPE usage_gas gauge
usage_gas 1019.003
```

## Development ##

Currently only the ESMR 5.0 format is supported and the parser is default configured to parse the telegram message with the keys the Sagemcom XS210 is using.
If you have to support a different ESMR 5.0 message, feel free to create your own implementation of the TelegramFormat struct. To support a different format then ESMR 5.0 you can implement your own implementation of the TelegramReaderOptions struct.
