package conn

import (
	"fmt"
	"io"
	"net/http"
)

type APISource struct {
	Endpoint string
}

func NewAPISource(endpoint string) Source {
	return &APISource{
		Endpoint: endpoint,
	}
}

func (s *APISource) ReadFromSource(telegramOptions *TelegramReaderOptions) (io.ReadCloser, error) {
	resp, err := http.Get(s.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("cannot read telegram from API endpoint - %w", err)
	}

	return resp.Body, nil
}
