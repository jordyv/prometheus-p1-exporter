package conn

import (
	"io"

	"github.com/jacobsa/go-serial/serial"
)

type SerialSource struct{}

func (SerialSource) ReadFromSource(telegramOptions *TelegramReaderOptions) (io.ReadWriteCloser, error) {
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

	return port, nil
}
