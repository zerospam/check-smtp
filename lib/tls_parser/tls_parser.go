package tls_parser

import (
	"crypto/tls"
	"fmt"
)

func DecodeString(str string) (uint16, error) {
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

func ToString(version uint16) (string, error) {
	switch version {
	case tls.VersionSSL30:
		return "VersionSSL30", nil
	case tls.VersionTLS10:
		return "VersionTLS10", nil
	case tls.VersionTLS11:
		return "VersionTLS11", nil
	case tls.VersionTLS12:
		return "VersionTLS12", nil
	default:
		return "", fmt.Errorf("%x unrecognized TLS version", version)
	}
}
