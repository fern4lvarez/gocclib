package cclib

import (
	"crypto/x509"
)

var (
	API_URL   = "https://api.devcctrl.com"
	SSL_CHECK = true
	CA_CERTS  *x509.CertPool
	CACHE     string // TODO
	DEBUG     = 0    // Set debug to 1 to enable debugging
	VERSION   = "0.0.5"
)
