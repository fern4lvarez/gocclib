package cclib

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
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
