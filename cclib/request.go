package cclib

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Request contains the API request basic information
type Request struct {
	email    string
	password string
	token    *Token
	version  string
	cache    string
	url      string
	sslCheck bool
	caCerts  *x509.CertPool
}

// New request creates a new api request having:
//
// * User email
//
// * User password
//
// * User token
//
// Returns a new request pointer
func NewRequest(email string, password string, token *Token) *Request {
	return &Request{
		email,
		password,
		token,
		VERSION,
		CACHE,
		API_URL,
		SSL_CHECK,
		CA_CERTS}
}

// Email returns request's email
func (request Request) Email() string {
	return request.email
}

// Password returns request's password
func (request Request) Password() string {
	return request.password
}

// Token returns request's token
func (request Request) Token() *Token {
	return request.token
}

// Cache returns request's cache
func (request Request) Cache() string {
	return request.cache
}

// Url returns request's Url
func (request Request) Url() string {
	return request.url
}

// SSLCheck returns true if SSL Certificate is checked and verified,
// returns false if SSL certificate check is skipped.
func (request Request) SSLCheck() bool {
	return request.sslCheck
}

// CaCerts returns request's root CA
func (request Request) CaCerts() *x509.CertPool {
	return request.caCerts
}

// SetEmail sets email address to a request
func (request *Request) SetEmail(email string) {
	request.email = email
}

// SetPassword sets a password to a request
func (request *Request) SetPassword(password string) {
	request.password = password
}

// SetToken sets a token to a request
func (request *Request) SetToken(token *Token) {
	request.token = token
}

// SetCache sets a cache to a request
func (request *Request) SetCache(cache string) {
	request.cache = cache
}

// SetUrl sets a URL to a request
func (request *Request) SetUrl(url string) {
	request.url = url
}

// EnableSSLCheck enables the SSL certificate verification
func (request *Request) EnableSSLCheck() {
	request.sslCheck = true
}

// DisableSSLCheck disables the SSL certificate verification
func (request *Request) DisableSSLCheck() {
	request.sslCheck = false
}

// SetCaCerts sets a set of root CA to a request
func (request *Request) SetCaCerts(caCerts *x509.CertPool) {
	request.caCerts = caCerts
}

// Post makes a POST request
func (request Request) Post(resource string, data url.Values) ([]byte, error) {
	return request.do(resource, "POST", data)
}

// Get makes a GET request
func (request Request) Get(resource string) ([]byte, error) {
	return request.do(resource, "GET", url.Values{})
}

// Put makes a PUT request
func (request Request) Put(resource string, data url.Values) ([]byte, error) {
	return request.do(resource, "PUT", data)
}

// Delete makes a DELETE request
func (request Request) Delete(resource string) ([]byte, error) {
	return request.do(resource, "DELETE", url.Values{})
}

func (request Request) do(resource string, method string, data url.Values) ([]byte, error) {
	u, err := url.ParseRequestURI(request.Url())
	if err != nil {
		return nil, err
	}
	u.Path = resource
	urlStr := fmt.Sprintf("%v", u)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !request.SSLCheck(),
			RootCAs:            request.CaCerts()},
	}
	client := &http.Client{Transport: tr}

	r, err := http.NewRequest(method, urlStr, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}

	if request.Token() != nil {
		r.Header.Add("Authorization", "cc_auth_token=\""+request.Token().Key()+"\"")
	} else if request.Email() != "" && request.Password() != "" {
		r.SetBasicAuth(request.Email(), request.Password())
	}
	r.Header.Add("Host", u.Host)
	r.Header.Add("User-Agent", "gocclib/"+Version())
	if m := strings.ToUpper(method); m == "POST" || m == "PUT" {
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	r.Header.Add("Accept-Encoding", "compress, gzip")

	if DEBUG {
		fmt.Printf("DEBUG Request >>> %v\n", r)
	}

	resp, err := client.Do(r)
	if err != nil {
		fmt.Printf("DEBUG Request Error >>> %v\n", err)
		return nil, err
	}

	if err = checkResponse(resp); err != nil {
		if DEBUG {
			fmt.Printf("DEBUG Request Error >>> %v\n", err)
			return nil, err
		}
	}

	defer resp.Body.Close()
	if DEBUG {
		fmt.Printf("DEBUG Response >>> %v\n", resp)
		fmt.Printf("DEBUG Body >>> %v\n", resp.Body)
	}

	return ioutil.ReadAll(resp.Body)
}
