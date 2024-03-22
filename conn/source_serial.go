package conn

import (
	"io"

	"github.com/jacobsa/go-serial/serial"
)

type SerialSource struct{
	UsbSerial string
}

func NewSerialSource(usbSerial string) Source {
	return &SerialSource{
		UsbSerial: usbSerial,
	}
}

func (s *SerialSource) ReadFromSource(telegramOptions *TelegramReaderOptions) (io.ReadCloser, error) {
	options := serial.OpenOptions{
		PortName:        s.UsbSerial,
		BaudRate:        telegramOptions.BaudRate,
		DataBits:        telegramOptions.DataBits,
		StopBits:        telegramOptions.StopBits,
		MinimumReadSize: 4,
	}

	port, err := serial.Open(options)
	if err != nil {
		return nil, err
	}

	return port, nil
}
