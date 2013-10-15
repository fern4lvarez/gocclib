package cclib

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"unicode/utf8"
)

var msgFail = "%v function fails. Expects %v, returns %v"

func decodeContent(content []byte) (data interface{}, err error) {
	if utf8.Valid(content) {
		json.Unmarshal(content, &data)
		return data, nil
	}

	buf := bytes.NewBuffer(content)
	reader, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(reader).Decode(&data)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return
}

func checkResponse(resp *http.Response) (err error) {
	switch resp.StatusCode {
	case 200, 201, 204:
		return nil
	default:
		return errors.New(resp.Status)
	}
}

func readerToStr(ir io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(ir)
	return buf.String()
}
