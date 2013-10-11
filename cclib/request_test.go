package cclib

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
	return c
}

func TestIt(t *testing.T) {
	b := mockHTTP([]byte(`{"token":"abcdefghijklmnopqrstuvxyz"}`))
	fmt.Println(b)
}
