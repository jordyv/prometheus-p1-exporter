package conn

import (
	"bufio"
	"io"
	"regexp"
)

type TelegramReaderOptions struct {
	BaudRate uint
	DataBits uint
	StopBits uint
}

type Source interface {
	ReadFromSource(telegramOptions *TelegramReaderOptions) (io.ReadWriteCloser, error)
}

func convertTelegramReaderToLines(telegramOptions *TelegramReaderOptions, reader io.ReadWriteCloser) ([]string, error) {
	buffer := bufio.NewReader(reader)
	telegram, err := buffer.ReadBytes('!')
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile("(\r)?\n")
	return re.Split(string(telegram), 100), err
}

func ReadTelegram(telegramOptions *TelegramReaderOptions, source Source) ([]string, error) {
	reader, err := source.ReadFromSource(telegramOptions)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return convertTelegramReaderToLines(telegramOptions, reader)
}
