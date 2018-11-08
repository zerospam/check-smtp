package mail_sender

import (
	"net"
)

type Conn struct {
	net.Conn
	firstRead     []byte
	firstLineRead bool
}

func (c *Conn) Read(b []byte) (n int, err error) {
	read, err := c.Conn.Read(b)
	if err != nil {
		return read, err
	}
	if c.firstLineRead {
		return read, err
	}

	c.firstRead = b
	c.firstLineRead = true

	return read, err
}

func NewConnection(conn net.Conn) *Conn {
	return &Conn{Conn: conn, firstLineRead: false}
}
