package cclib

import (
	"testing"
)

func TestVersion(t *testing.T) {
	// When
	v := Version()

	// Then
	if v != "0.4.0" {
		t.Errorf(msgFail, "Version", "0.4.0", v)
	}
}
