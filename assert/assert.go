package assert

import "testing"

func AssertEq(t *testing.T, a, b interface{}) {
	if a != b {
		t.Fatalf("Assertion error: %+v != %+v", a, b)
	}
}

func AssertNotEq(t *testing.T, a, b interface{}) {
	if a == b {
		t.Fatalf("Assertion error: %+v == %+v", a, b)
	}
}
