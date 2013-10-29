package cclib

import (
	"testing"
)

func TestVersion(t *testing.T) {
	// When
	v := Version()

	// Then
	if v != "0.0.5" {
		t.Errorf(msgFail, "Version", "0.0.5", v)
	}
}
