package conn

import (
	"testing"
)

func TestReadTelegram(t *testing.T) {
	lines, err := ReadTelegram(&Esmr5TelegramReaderOptions, &MockSource{})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(lines) != 27 {
		t.Fatal("expect 27 lines in telegram body")
	}
}
