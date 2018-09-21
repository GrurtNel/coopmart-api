package utils

import (
	"testing"
)

func TestUnixToDate(t *testing.T) {
	if UnixToDate(1516731454) != "24-01-2018 01:17:34" {
		t.Log(UnixToDate(1516731454))
		t.Fail()
	} else {
		t.Log(UnixToDate(1516731454))
	}
}
