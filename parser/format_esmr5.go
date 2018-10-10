package parser

var Esmr5TelegramFormat = TelegramFormat{
	LineTimestamp:                   4,
	LineElectricityUsageLow:         6,
	LineElectricityUsageHigh:        7,
	LineElectricityReturnedLow:      8,
	LineElectricityReturnedHigh:     9,
	LineActiveTariff:                10,
	LineCurrentElectricityUsageLow:  11,
	LineCurrentElectricityUsageHigh: 12,
	LinePowerFailuresShort:          13,
	LinePowerFailuresLong:           14,
	LineGasUsage:                    25,
}
