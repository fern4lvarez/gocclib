package cclib

import (
	"testing"
)

func TestVersion(t *testing.T) {
	// When
	v := Version()

	// Then
	if v != "0.1.0" {
		t.Errorf(msgFail, "Version", "0.1.0", v)
	}
}
