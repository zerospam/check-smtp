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
	conn, err = net.DialTimeout("tcp", t.Address(), timeout)
	return conn, err
}
