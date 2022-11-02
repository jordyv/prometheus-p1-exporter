package conn

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strings"
)

type TelegramReaderOptions struct {
	BaudRate uint
	DataBits uint
	StopBits uint
}

type Source interface {
	ReadFromSource(telegramOptions *TelegramReaderOptions) (io.ReadCloser, error)
}

func convertTelegramReaderToLines(telegramOptions *TelegramReaderOptions, reader io.ReadCloser) ([]string, error) {
	buffer := bufio.NewReader(reader)
	telegram, err := buffer.ReadBytes('!')
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile("(\r)?\n")
	return re.Split(string(telegram), 100), err
}

func splitTelegramLine(line string) (string, string, error) {
	re := regexp.MustCompile("^([0-9\\-\\:\\.]*)(\\(.*\\))*(\\(.*\\))$")
	lineParts := re.FindStringSubmatch(line)
	if len(lineParts) < 4 {
		return "", "", errors.New("could not split telegram line correctly")
	}
	key := lineParts[1]
	value := strings.TrimPrefix(strings.TrimSuffix(lineParts[len(lineParts)-1], ")"), "(")

	return key, value, nil
}

func buildTelegramMap(lines []string) map[string]string {
	telegramMap := make(map[string]string, 0)
	for _, line := range lines {
		key, value, err := splitTelegramLine(line)
		if err == nil {
			telegramMap[key] = value
		}
	}
	return telegramMap
}

func ReadTelegram(telegramOptions *TelegramReaderOptions, source Source) (map[string]string, error) {
	reader, err := source.ReadFromSource(telegramOptions)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	lines, err := convertTelegramReaderToLines(telegramOptions, reader)
	if err != nil {
		return nil, err
	}
	return buildTelegramMap(lines), nil
}
