package cclib

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"unicode/utf8"
)

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
