package cclib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Request struct {
	email           string
	password        string
	token           *Token
	version         string
	cache           string
	url             string
	disableSSLCheck bool
	caCerts         string
}

func NewRequest(email string, password string, token *Token) *Request {
	return &Request{
		email,
		password,
		token,
		VERSION,
		CACHE,
		API_URL,
		DISABLE_SSL_CHECK,
		CA_CERTS}
}

func (request Request) Email() string {
	return request.email
}

func (request Request) Password() string {
	return request.password
}

func (request Request) Token() *Token {
	return request.token
}

func (request Request) Cache() string {
	return request.cache
}

func (request Request) Url() string {
	return request.url
}

func (request Request) DisableSSLCheck() bool {
	return request.disableSSLCheck
}

func (request Request) CaCerts() string {
	return request.caCerts
}

func (request Request) Post(resource string, data url.Values) ([]byte, error) {
	return request.do(resource, "POST", data)
}

func (request Request) Get(resource string) ([]byte, error) {
	return request.do(resource, "GET", url.Values{})
}

func (request Request) Put(resource string, data url.Values) ([]byte, error) {
	return request.do(resource, "PUT", data)
}

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
	client := &http.Client{}

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

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	if err = checkResponse(resp); err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
