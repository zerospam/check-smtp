package lib

import (
	"crypto/tls"
	"github.com/zerospam/check-firewall/lib/tls-generator"
	"time"
)

type TransportServer struct {
	Server       string `json:"server"`
	Port         int    `json:"port"`
	OnMx         bool   `json:"mx"`
	TestEmail    string `json:"test_email"`
	tlsGenerator *tlsgenerator.CertificateGenerator
}

func (t *TransportServer) getClientTLSConfig(commonName string) *tls.Config {
	if t.tlsGenerator == nil {
		t.tlsGenerator = tlsgenerator.NewClient(time.Now(), 30*365*24*time.Hour)
	}

	return t.tlsGenerator.GetTlsClientConfig(commonName)
}
