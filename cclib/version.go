package cclib

import (
	"io/ioutil"
)

// Version of the library
func Version() string {
	b, _ := ioutil.ReadFile("VERSION")
	v := string(b)
	return v
}
