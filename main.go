package main

import (
	"bufio"
	"fmt"
	"log"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

// TODO: Increase this interval
const ReadInterval = 2 * time.Second

/*
Example:
<START>
1
2	/Ene5\XS210 ESMR 5.0
3
4	1-3:0.2.8(50)													<- header
5	0-0:1.0.0(181009214805S)										<- Timestamp YYMMDDhhmmssX
6	0-0:96.1.1(xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx)					<- equipment ID
7	1-0:1.8.1(001179.186*kWh)										<- electricity used tariff 1
8	1-0:1.8.2(001225.590*kWh)										<- electricity used tariff 2
9	1-0:2.8.1(000000.016*kWh)										<- electricity delivered tariff 1
10	1-0:2.8.2(000000.000*kWh)										<- electricity delivered tariff 2
11	0-0:96.14.0(0002)												<- active tariff
12	1-0:1.7.0(00.200*kW)											<- electricity current usage
13	1-0:2.7.0(00.000*kW)											<- electricity current delivery
14	0-0:96.7.21(00057)												<- short power failure count
15	0-0:96.7.9(00002)												<- long power failure count
16	1-0:99.97.0(1)(0-0:96.7.19)(170829233732S)(0000001803*s)		<- ?
17	1-0:32.32.0(00002)												<- voltage sag l1 count
18	1-0:32.36.0(00000)												<- voltage swell l1 count
19	0-0:96.13.0()													<- text message
20	1-0:32.7.0(227.0*V)												<- instantaneous voltage
21	1-0:31.7.0(001*A)												<- instantaneous current l1
22	1-0:21.7.0(00.200*kW)											<- instantaneous active power l1 positive
23	1-0:22.7.0(00.000*kW)											<- instantaneous active power l1 negative
24	0-1:24.1.0(003)													<- device type
25	0-1:96.1.0(xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx)					<- identifier gas
26	0-1:24.2.1(181009214500S)(01019.003*m3)							<- hourly gas meter reading
27	!6611
<END>
*/

type TelegramReaderOptions struct {
	BaudRate     uint
	DataBits     uint
	StopBits     uint
	TelegramSize int
}

type TelegramFormat struct {
	LineTimestamp       int
	LineElectricityLow  int
	LineElectricityHigh int
	LineGas             int
}

var Esmr5TelegramFormat = TelegramFormat{LineTimestamp: 5, LineElectricityHigh: 7, LineElectricityLow: 8, LineGas: 26}

type Telegram struct {
	Timestamp       string
	ElectricityLow  string
	ElectricityHigh string
	Gas             string
}

var Esmr5TelegramReaderOptions = TelegramReaderOptions{BaudRate: 115200, DataBits: 8, StopBits: 1, TelegramSize: 27}

func readTelegram(telegramOptions *TelegramReaderOptions) ([]string, error) {
	options := serial.OpenOptions{
		PortName:        "/dev/ttyUSB0",
		BaudRate:        telegramOptions.BaudRate,
		DataBits:        telegramOptions.DataBits,
		StopBits:        telegramOptions.StopBits,
		MinimumReadSize: 4,
	}

	port, err := serial.Open(options)
	if err != nil {
		return nil, err
	}

	defer port.Close()

	reader := bufio.NewReader(port)
	linesToRead := telegramOptions.TelegramSize
	lines := make([]string, telegramOptions.TelegramSize)
	for i := 1; i <= linesToRead; i++ {
		line, _, _ := reader.ReadLine()
		lines = append(lines, string(line))
	}
	return lines, nil
}

func parseTelegram(format *TelegramFormat, telegramLines []string) Telegram {
	return Telegram{
		Timestamp:       telegramLines[format.LineTimestamp],
		ElectricityHigh: telegramLines[format.LineElectricityHigh],
		ElectricityLow:  telegramLines[format.LineElectricityLow],
		Gas:             telegramLines[format.LineGas],
	}
}

func main() {
	for {
		lines, err := readTelegram(&Esmr5TelegramReaderOptions)
		if err != nil {
			log.Fatalln("Error while opening serial device: " + err.Error())
			continue
		}
		fmt.Println(parseTelegram(&Esmr5TelegramFormat, lines))
		// TODO: Parse telegram to struct
		time.Sleep(ReadInterval)
	}
}
