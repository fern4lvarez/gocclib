package cclib

import (
	"testing"
)

func TestVersion(t *testing.T) {
	v := Version()
	if v != "0.0.1" {
		t.Errorf(msgFail, "Version", "0.0.1", v)
	}
}
