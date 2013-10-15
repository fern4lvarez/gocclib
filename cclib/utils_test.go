package cclib

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func TestDecodeContentUTF8Valid(t *testing.T) {
	// Given
	content := []byte(`{
					"foo":"abcdefghijklmnopqrstuvxyz",
				  	"bar":"1234567890"
				  }`)

	var expectedData interface{}
	err1 := json.Unmarshal(content, &expectedData)

	// When
	data, err2 := decodeContent(content)

	// Then
	if err1 != nil {
		t.Errorf(msgFail, "decodeContent", nil, err1)
	}
	if err2 != nil {
		t.Errorf(msgFail, "decodeContent", nil, err2)
	}
	if !reflect.DeepEqual(data, expectedData) {
		t.Errorf(msgFail, "decodeContent", expectedData, data)
	}
}

func TestDecodeContentUTF8Invalid(t *testing.T) {
	// Given
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write([]byte(`{
					"foo":"abcdefghijklmnopqrstuvxyz",
				  	"bar":"1234567890"
				  }`))
	w.Close()
	content := b.Bytes() // gzip content

	c := []byte(`{
					"foo":"abcdefghijklmnopqrstuvxyz",
				  	"bar":"1234567890"
				  }`)

	var expectedData interface{}
	err1 := json.Unmarshal(c, &expectedData)

	// When
	data, err2 := decodeContent(content)

	// Then
	if err1 != nil {
		t.Errorf(msgFail, "decodeContent", nil, err1)
	}
	if err2 != nil {
		t.Errorf(msgFail, "decodeContent", nil, err2)
	}
	if !reflect.DeepEqual(data, expectedData) {
		t.Errorf(msgFail, "decodeContent", expectedData, data)
	}
}

func TestCheckResponse(t *testing.T) {
	// Given
	resp200 := &http.Response{StatusCode: 200, Status: "200 OK"}
	resp201 := &http.Response{StatusCode: 201, Status: "201 Created"}
	resp204 := &http.Response{StatusCode: 204, Status: "204 No Content"}
	resp404 := &http.Response{StatusCode: 404, Status: "404 Not Found"}

	// When
	err200 := checkResponse(resp200)
	err201 := checkResponse(resp201)
	err204 := checkResponse(resp204)
	err404 := checkResponse(resp404)

	// Then
	if err200 != nil {
		t.Errorf(msgFail, "checkResponse", nil, err200)
	}
	if err201 != nil {
		t.Errorf(msgFail, "checkResponse", nil, err201)
	}
	if err204 != nil {
		t.Errorf(msgFail, "checkResponse", nil, err204)
	}
	if err404 == nil {
		t.Errorf(msgFail, "checkResponse", errors.New(resp404.Status), err404)
	}
}

func TestReaderToStr(t *testing.T) {
	// Given
	s := "Hello"
	ir := bytes.NewBufferString(s)

	// When
	rts := readerToStr(ir)

	// Then
	if rts != s {
		t.Errorf(msgFail, "readerToStr", s, rts)
	}
}
