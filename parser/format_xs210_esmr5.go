package parser

var XS210ESMR5TelegramFormat = TelegramFormat{
	KeyTimestamp:                   "0-0:1.0.0",
	KeyElectricityUsageHigh:        "1-0:1.8.1",
	KeyElectricityUsageLow:         "1-0:1.8.2",
	KeyElectricityReturnedHigh:     "1-0:2.8.1",
	KeyElectricityReturnedLow:      "1-0:2.8.2",
	KeyActiveTariff:                "0-0:96.14.0",
	KeyCurrentElectricityUsageHigh: "1-0:1.7.0",
	KeyCurrentElectricityUsageLow:  "1-0:2.7.0",
	KeyPowerFailuresShort:          "0-0:96.7.21",
	KeyPowerFailuresLong:           "0-0:96.7.9",
	KeyGasUsage:                    "0-1:24.2.1",
}
