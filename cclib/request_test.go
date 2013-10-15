package cclib

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequest(t *testing.T) {
	// Given
	email := "user@example.com"
	password := "password"
	token := &Token{
		"token": "1234567890",
	}

	// When
	req := NewRequest(email, password, token)

	// Then
	if req.Email() != email {
		t.Errorf(msgFail, "NewRequest and Email", email, req.Email())
	}
	if req.Password() != password {
		t.Errorf(msgFail, "NewRequest and Password", password, req.Password())
	}
	if req.Token() != token {
		t.Errorf(msgFail, "NewRequest and Token", token, req.Token())
	}
	if req.Cache() != CACHE {
		t.Errorf(msgFail, "NewRequest and Cache", CACHE, req.Cache())
	}
	if req.Url() != API_URL {
		t.Errorf(msgFail, "NewRequest and Url", API_URL, req.Url())
	}
	if req.DisableSSLCheck() != DISABLE_SSL_CHECK {
		t.Errorf(msgFail, "NewRequest and DisableSSLCheck", DISABLE_SSL_CHECK, req.DisableSSLCheck())
	}
	if req.CaCerts() != CA_CERTS {
		t.Errorf(msgFail, "NewRequest and CaCerts", CA_CERTS, req.CaCerts())
	}
}

func mockHTTP(content []byte) []byte {
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, content)
	}

	req, err := http.NewRequest("GET", API_URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler(w, req)

	c, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(c))
	return c
}
