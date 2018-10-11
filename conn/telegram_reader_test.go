package conn

import (
	"testing"
)

func TestReadTelegram(t *testing.T) {
	items, err := ReadTelegram(&ESMR5TelegramReaderOptions, &MockSource{})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(items) != 23 {
		t.Fatal("expect 23 items in telegram body map")
	}
}
