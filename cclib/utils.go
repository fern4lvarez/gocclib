package cclib

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"
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

func readCredentialsFile(filepath string) (email, password string, err error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", "", err
	}

	scanner := bufio.NewScanner(f)
	lines := 0
	for scanner.Scan() {
		switch lines {
		case 0:
			email = scanner.Text()
			lines++
		case 1:
			password = scanner.Text()
			lines++
		case 2:
			break
		}
	}

	if lines != 2 {
		return "", "", errors.New("Not lines enough on credentials file.")
	}

	if err := scanner.Err(); err != nil {
		return "", "", err
	}

	return
}

func buildTimestamp(dt *time.Time) string {
	m := strconv.Itoa(dt.Nanosecond())
	if len(m) > 5 {
		m = m[:6]
	}

	u := dt.Unix()

	return fmt.Sprintf("%s.%s", strconv.Itoa(int(u)), m)
}

// taken from http://play.golang.org/p/fpkK48W9Rp
func isNil(value interface{}) bool {
	if value == nil {
		return true
	}
	if !reflect.ValueOf(value).Elem().IsValid() {
		return true
	}

	return false
}
