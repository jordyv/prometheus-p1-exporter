package parser

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
)

type Telegram struct {
	Timestamp                   int64
	ElectricityUsageLow         float64
	ElectricityUsageHigh        float64
	ElectricityReturnedLow      float64
	ElectricityReturnedHigh     float64
	ActiveTariff                int64
	PowerFailuresShort          int64
	PowerFailuresLong           int64
	CurrentElectricityUsageLow  float64
	CurrentElectricityUsageHigh float64
	GasUsage                    float64
}

type TelegramFormat struct {
	LineTimestamp                   int
	LineElectricityUsageLow         int
	LineElectricityUsageHigh        int
	LineElectricityReturnedLow      int
	LineElectricityReturnedHigh     int
	LineActiveTariff                int
	LinePowerFailuresShort          int
	LinePowerFailuresLong           int
	LineCurrentElectricityUsageLow  int
	LineCurrentElectricityUsageHigh int
	LineGasUsage                    int
}

func parseTelegramValue(s string) string {
	re := regexp.MustCompile("^\\d+\\-\\d+\\:\\d+\\.\\d+\\.\\d+(\\(.*\\))?\\((.*)\\)$")
	parsed := re.FindStringSubmatch(s)
	if len(parsed) > 0 {
		return parsed[len(parsed)-1]
	}
	return s
}

func parseInt(s string) int64 {
	res, _ := strconv.ParseInt(strings.TrimLeft(parseTelegramValue(s), "0"), 0, 64)
	return res
}

func parseTimestampString(s string) int64 {
	// 0-0:1.0.0(181009214805S)
	res, _ := strconv.ParseInt(strings.Replace(parseTelegramValue(s), "S", "", 1), 0, 64)
	return res
}

func parseElectricityStringWithSuffix(s string, suffix string) float64 {
	res, _ := strconv.ParseFloat(strings.Replace(parseTelegramValue(s), suffix, "", 1), 64)
	return res
}

func parseElectricityString(s string) float64 {
	// 1-0:1.8.1(001179.186*kWh)
	// 1-0:1.8.2(001225.590*kWh)
	return parseElectricityStringWithSuffix(s, "*kWh")
}

func parseGasString(s string) float64 {
	// 0-1:24.2.1(181009214500S)(01019.003*m3)
	res, _ := strconv.ParseFloat(strings.Replace(parseTelegramValue(s), "*m3", "", 1), 64)
	return res
}

func ParseTelegram(format *TelegramFormat, telegramLines []string) (Telegram, error) {
	logrus.Debugln("Line to parse", telegramLines)

	if len(telegramLines) > 0 {
		return Telegram{
			Timestamp:                   parseTimestampString(telegramLines[format.LineTimestamp]),
			ElectricityUsageHigh:        parseElectricityString(telegramLines[format.LineElectricityUsageHigh]),
			ElectricityUsageLow:         parseElectricityString(telegramLines[format.LineElectricityUsageLow]),
			ElectricityReturnedHigh:     parseElectricityString(telegramLines[format.LineElectricityReturnedHigh]),
			ElectricityReturnedLow:      parseElectricityString(telegramLines[format.LineElectricityReturnedLow]),
			ActiveTariff:                parseInt(telegramLines[format.LineActiveTariff]),
			PowerFailuresLong:           parseInt(telegramLines[format.LinePowerFailuresLong]),
			PowerFailuresShort:          parseInt(telegramLines[format.LinePowerFailuresShort]),
			CurrentElectricityUsageHigh: parseElectricityStringWithSuffix(telegramLines[format.LineCurrentElectricityUsageHigh], "*kW"),
			CurrentElectricityUsageLow:  parseElectricityStringWithSuffix(telegramLines[format.LineCurrentElectricityUsageLow], "*kW"),
			GasUsage:                    parseGasString(telegramLines[format.LineGasUsage]),
		}, nil
	}
	return Telegram{}, errors.New("provided telegram is empty")
}
