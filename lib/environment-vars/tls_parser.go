package environmentvars

import (
	"crypto/tls"
	"fmt"
)

func tlsDecodeString(str string) (uint16, error) {
	switch str {
	case "SSL30":
		return tls.VersionSSL30, nil
	case "TLS10":
		return tls.VersionTLS10, nil
	case "TLS11":
		return tls.VersionTLS11, nil
	case "TLS12":
		return tls.VersionTLS12, nil
	default:
		return 0, fmt.Errorf("%s unrecognized TLS version", str)
	}
}
