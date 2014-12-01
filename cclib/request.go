package cclib

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Request contains the API request basic information
type Request struct {
	Email    string
	Password string
	SslCheck bool
	Api      Api
	CaCerts  *x509.CertPool
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
func NewRequest(email string, password string, api Api) *Request {
	return &Request{
		email,
		password,
		SSL_CHECK,
		api,
		CA_CERTS}
}

// SetEmail sets email address to a request
func (request *Request) SetEmail(email string) {
	request.Email = email
}

// SetPassword sets a password to a request
func (request *Request) SetPassword(password string) {
	request.Password = password
}

// EnableSSLCheck enables the SSL certificate verification
func (request *Request) EnableSSLCheck() {
	request.SslCheck = true
}

// DisableSSLCheck disables the SSL certificate verification
func (request *Request) DisableSSLCheck() {
	request.SslCheck = false
}

// SetApi sets an Api to a request
func (request *Request) SetApi(api Api) {
	request.Api = api
}

// SetCaCerts sets a set of root CA to a request
func (request *Request) SetCaCerts(caCerts *x509.CertPool) {
	request.CaCerts = caCerts
}

// Post makes a POST request
func (request Request) Post(resource string, data url.Values) ([]byte, error) {
	return request.do(resource, "POST", []byte(data.Encode()), false, false)
}

// Get makes a GET request
func (request Request) Get(resource string) ([]byte, error) {
	return request.do(resource, "GET", []byte{}, false, false)
}

// Put makes a PUT request
func (request Request) Put(resource string, data url.Values) ([]byte, error) {
	return request.do(resource, "PUT", []byte(data.Encode()), false, false)
}

// Delete makes a DELETE request
func (request Request) Delete(resource string) ([]byte, error) {
	return request.do(resource, "DELETE", []byte{}, false, false)
}

// HeadToken makes an auth HEAD request to a regular
// endpoint to check if token is still valid
func (request Request) HeadToken() error {
	_, err := request.do("/user/", "HEAD", nil, false, false)
	return err
}

// PostToken makes a POST request to the token source URL
func (request Request) PostToken() ([]byte, error) {
	return request.do("", "POST", nil, true, false)
}

// PostToken makes a POST request to the register add-on URL
func (request Request) PostAddon(data []byte) ([]byte, error) {
	return request.do("", "POST", data, false, true)
}

func (request Request) doUrl(isTokenReq bool, isAddonReq bool) string {
	switch {
	case isTokenReq:
		return request.Api.TokenSourceUrl()
	case isAddonReq:
		return request.Api.RegisterAddonUrl()
	}
	return request.Api.Url()
}

func (request Request) do(resource string, method string, data []byte, isTokenReq bool, isAddonReq bool) ([]byte, error) {
	request_url := request.doUrl(isTokenReq, isAddonReq)
	u, err := url.ParseRequestURI(request_url)
	if err != nil {
		return nil, err
	}

	if resource != "" {
		u.Path = resource
	}

	urlStr := fmt.Sprintf("%v", u)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !request.SslCheck,
			RootCAs:            request.CaCerts},
	}
	client := &http.Client{Transport: tr}

	r, err := http.NewRequest(method, urlStr, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	if !isNil(request.Api.Token()) {
		r.Header.Add("Authorization", "cc_auth_token=\""+request.Api.Token().Key+"\"")
	} else if request.Email != "" && request.Password != "" {
		r.SetBasicAuth(request.Email, request.Password)
	} else {
		return nil, errors.New("Request not authorized.")
	}

	r.Header.Add("Host", u.Host)
	r.Header.Add("User-Agent", "gocclib/"+Version())
	if m := strings.ToUpper(method); m == "POST" || m == "PUT" {
		contentType := "application/x-www-form-urlencoded"
		if isJSONData(data) {
			contentType = "application/json"
		}
		r.Header.Add("Content-Type", contentType)
	}
	r.Header.Add("Content-Length", strconv.Itoa(len(data)))
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
		}
		return nil, err
	}

	defer resp.Body.Close()
	if DEBUG {
		fmt.Printf("DEBUG Response >>> %v\n", resp)
		fmt.Printf("DEBUG Body >>> %v\n", resp.Body)
	}

	return ioutil.ReadAll(resp.Body)
}
