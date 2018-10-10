package conn

import (
	"io"
	"strings"
)

var dummyTelegram = `
/Ene5\XS210 ESMR 5.0

1-3:0.2.8(50)
0-0:1.0.0(181009214805S)
0-0:96.1.1(xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx)
1-0:1.8.1(001179.186*kWh)
1-0:1.8.2(001225.590*kWh)
1-0:2.8.1(000000.016*kWh)
1-0:2.8.2(000000.000*kWh)
0-0:96.14.0(0002)
1-0:1.7.0(00.200*kW)
1-0:2.7.0(00.000*kW)
0-0:96.7.21(00057)
0-0:96.7.9(00002)
1-0:99.97.0(1)(0-0:96.7.19)(170829233732S)(0000001803*s)
1-0:32.32.0(00002)
1-0:32.36.0(00000)
0-0:96.13.0()
1-0:32.7.0(227.0*V)
1-0:31.7.0(001*A)
1-0:21.7.0(00.200*kW)
1-0:22.7.0(00.000*kW)
0-1:24.1.0(003)
0-1:96.1.0(xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx)
0-1:24.2.1(181009214500S)(01019.003*m3)
!6611
`

type MockSource struct{}

type MockSourceReader struct {
	Content string
}

func (r MockSourceReader) Read(b []byte) (int, error) {
	return strings.NewReader(r.Content).Read(b)
}
func (r MockSourceReader) Write(b []byte) (int, error) {
	return 0, nil
}
func (MockSourceReader) Close() error {
	return nil
}

func (MockSource) ReadFromSource(telegramOptions *TelegramReaderOptions) (io.ReadWriteCloser, error) {
	return MockSourceReader{Content: dummyTelegram}, nil
}
