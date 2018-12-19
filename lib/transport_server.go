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

func (t *TransportServer) Address() (string, error) {

	if !t.OnMx {
		return fmt.Sprintf("%s:%d", t.Server, t.Port), nil
	}

	mxRecords, errorMx := net.LookupMX(t.Server)
	if errorMx != nil {
		return "", errorMx
	}
	if len(mxRecords) == 0 {
		return "", fmt.Errorf("MX Records: no mx records found for %s", t.Server)
	}
	return fmt.Sprintf("%s:%d", mxRecords[0].Host, t.Port), nil

}

func (t *TransportServer) Connect(timeout time.Duration) (conn net.Conn, err error) {
	address, err := t.Address()
	if err != nil {
		return nil, err
	}

	return net.DialTimeout("tcp", address, timeout)
}
