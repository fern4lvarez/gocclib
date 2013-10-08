package cclib

import (
	"encoding/json"
)

type Token map[string]string

func (t *Token) Decode(b []byte) (err error) {
	return json.Unmarshal(b, &t)
}

func (t Token) Encode() (b []byte, err error) {
	b, err = json.Marshal(t)
	return
}

func (t Token) Key() string {
	return t["token"]
}
