package cclib

import (
	"crypto/x509"
)

var (
	API_URL   = "https://api.cloudcontrol.com"
	SSL_CHECK = true
	CA_CERTS  *x509.CertPool
	CACHE     string  // TODO
	DEBUG     = false // Set debug to true to enable debugging
	VERSION   = "0.2.1"
)
