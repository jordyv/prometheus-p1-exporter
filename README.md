# Prometheus P1 exporter

Prometheus exporter for smart meter statistics fetched with a P1 cable.


## Installation

### From source

With Go get:

```shell
go install github.com/loafoe/prometheus-p1-exporter@latest
```

## Usage

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

By default, the exporter will collect metrics from `/dev/ttyUSB0` every 10 seconds and export the metrics to an HTTP endpoint at `http://127.0.0.1:8888/metrics`. This endpoint can be added to your Prometheus configuration.

Example metrics page:

```
# HELP p1_active_tariff Active tariff
# TYPE p1_active_tariff gauge
p1_active_tariff 2
# HELP p1_current_usage_electricity_high Electricity currently used high tariff
# TYPE p1_current_usage_electricity_high gauge
p1_current_usage_electricity_high 0
# HELP p1_current_usage_electricity_low Electricity currently used low tariff
# TYPE p1_current_usage_electricity_low gauge
p1_current_usage_electricity_low 0.2
# HELP p1_power_failures_long Power failures long
# TYPE p1_power_failures_long gauge
p1_power_failures_long 2
# HELP p1_power_failures_short Power failures short
# TYPE p1_power_failures_short gauge
p1_power_failures_short 57
# HELP p1_returned_electricity_high Electricity returned high tariff
# TYPE p1_returned_electricity_high gauge
p1_returned_electricity_high 0
# HELP p1_returned_electricity_low Electricity returned low tariff
# TYPE p1_returned_electricity_low gauge
p1_returned_electricity_low 0.016
# HELP p1_usage_electricity_high Electricity usage high tariff
# TYPE p1_usage_electricity_high gauge
p1_usage_electricity_high 1225.59
# HELP p1_usage_electricity_low Electricity usage low tariff
# TYPE p1_usage_electricity_low gauge
p1_usage_electricity_low 1179.186
# HELP p1_usage_gas Gas usage
# TYPE p1_usage_gas gauge
p1_usage_gas 1019.003
```

## Development

Currently only the ESMR 5.0 format is supported and the parser is default configured to parse the telegram message with the keys the Sagemcom XS210 is using.
If you have to support a different ESMR 5.0 message, feel free to create your own implementation of the TelegramFormat struct. To support a different format then ESMR 5.0 you can implement your own implementation of the TelegramReaderOptions struct.

## Acknowledgement

This was forked from https://github.com/jordyv/prometheus-p1-exporter
