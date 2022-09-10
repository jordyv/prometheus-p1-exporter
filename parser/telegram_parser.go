package parser

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type Telegram struct {
	Timestamp               int64
	ElectricityUsageLow     float64
	ElectricityUsageHigh    float64
	ElectricityReturnedLow  float64
	ElectricityReturnedHigh float64
	ActiveTariff            int64
	PowerFailuresShort      int64
	PowerFailuresLong       int64
	ActiveElectricityDraw   float64
	ActiveElectricityReturn float64
	GasUsage                float64
}

type TelegramFormat struct {
	KeyTimestamp               string
	KeyElectricityUsageLow     string
	KeyElectricityUsageHigh    string
	KeyElectricityReturnedLow  string
	KeyElectricityReturnedHigh string
	KeyActiveTariff            string
	KeyActiveElectricityDraw   string
	KeyActiveElectricityReturn string
	KeyPowerFailuresShort      string
	KeyPowerFailuresLong       string
	KeyGasUsage                string
}

func parseTelegramValue(s string) string {
	re := regexp.MustCompile(`^\d+-\d+:\d+\.\d+\.\d+(\(.*\))?\((.*)\)$`)
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

func ParseTelegram(format *TelegramFormat, telegramLines map[string]string) (Telegram, error) {
	logrus.Debugln("Line to parse", telegramLines)

	if len(telegramLines) > 0 {
		return Telegram{
			Timestamp:               parseTimestampString(telegramLines[format.KeyTimestamp]),
			ElectricityUsageHigh:    parseElectricityString(telegramLines[format.KeyElectricityUsageHigh]),
			ElectricityUsageLow:     parseElectricityString(telegramLines[format.KeyElectricityUsageLow]),
			ElectricityReturnedHigh: parseElectricityString(telegramLines[format.KeyElectricityReturnedHigh]),
			ElectricityReturnedLow:  parseElectricityString(telegramLines[format.KeyElectricityReturnedLow]),
			ActiveTariff:            parseInt(telegramLines[format.KeyActiveTariff]),
			PowerFailuresLong:       parseInt(telegramLines[format.KeyPowerFailuresLong]),
			PowerFailuresShort:      parseInt(telegramLines[format.KeyPowerFailuresShort]),
			ActiveElectricityDraw:   parseElectricityStringWithSuffix(telegramLines[format.KeyActiveElectricityDraw], "*kW"),
			ActiveElectricityReturn: parseElectricityStringWithSuffix(telegramLines[format.KeyActiveElectricityReturn], "*kW"),
			GasUsage:                parseGasString(telegramLines[format.KeyGasUsage]),
		}, nil
	}
	return Telegram{}, errors.New("provided telegram is empty")
}
