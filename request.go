package cclib

import (
	"bytes"
	"io"
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

func (request Request) Req(resource string, method string, data io.Reader) ([]byte, error) {
	urlStr := request.Url() + resource
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	r, err := http.NewRequest(method, urlStr, data)
	if err != nil {
		return nil, err
	}

	if request.Token() != nil {
		r.Header.Add("Authorization", "cc_auth_token=\""+request.Token().Key()+"\"")
	} else if request.Email() != "" && request.Password() != "" {
		r.SetBasicAuth(request.Email(), request.Password())
	}
	r.Header.Add("Host", u.Host)
	r.Header.Add("User-Agent", "gocclib/0.0.1")
	if m := strings.ToUpper(method); m == "POST" || m == "PUT" {
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	r.Header.Add("Content-Length", strconv.Itoa(len(readerToStr(data))))
	r.Header.Add("Accept-Encoding", "compress, gzip")

	resp, err := client.Do(r)
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
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

func (request Request) Post(resource string, data []byte) ([]byte, error) {
	if data == nil {
		data = []byte{}
	}
	return request.Req(resource, "POST", bytes.NewReader(data))
}

func (request Request) Get(resource string) ([]byte, error) {
	return request.Req(resource, "GET", bytes.NewReader([]byte{}))
}

func (request Request) Put(resource string, data []byte) ([]byte, error) {
	if data == nil {
		data = []byte{}
	}
	return request.Req(resource, "PUT", bytes.NewReader(data))
}

func readerToStr(ir io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(ir)
	return buf.String()
}
