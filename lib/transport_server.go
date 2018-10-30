package lib

import (
	"fmt"
	"net"
	"time"
)

type TransportServer struct {
	Server    string `json:"server"`
	Port      int    `json:"port"`
	OnMx      bool   `json:"mx"`
	TestEmail string `json:"test_email"`
}

func (t *TransportServer) Address() string {
	return fmt.Sprintf("%s:%d", t.Server, t.Port)
}

func (t *TransportServer) Connect(timeout time.Duration) (conn net.Conn, err error) {
	address := t.Address()

	if t.OnMx {
		mxRecords, errorMx := net.LookupMX(t.Server)
		if errorMx != nil {
			return nil, errorMx
		}
		if len(mxRecords) == 0 {
			return nil, fmt.Errorf("MX Records: no mx records found for %s", t.Server)
		}
		address = mxRecords[0].Host
	}

	conn, err = net.DialTimeout("tcp", address, timeout)
	return conn, err
}
