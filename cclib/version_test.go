package cclib

import (
	"testing"
)

func TestVersion(t *testing.T) {
	// When
	v := Version()

	// Then
	if v != "0.2.1" {
		t.Errorf(msgFail, "Version", "0.2.1", v)
	}
}
