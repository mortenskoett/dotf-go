// Contains helpers used for testing
package test

import (
	"reflect"
	"testing"
)

func AssertEqual(want, got interface{}, t *testing.T) {
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected: %v, got: %v", want, got)
	}
}

// Fail fails hard, prints actual and expected and panics. Good to stop flow in a test with multiple
// assertions.
func FailHard(actual, expected interface{}, t *testing.T) {
	t.Fatalf("\nactual = %+v\nexpected = %+v", actual, expected)
}
// Fail fails hard, prints message, actual and expected and panics. Good to stop flow in a test with
// multiple assertions.
func FailHardMsg(msg string, actual, expected interface{}, t *testing.T) {
	t.Fatalf("%s: \nactual = %+v\nexpected = %+v", msg, actual, expected)
}

// Fail fails and prints actual and expected
func Fail(actual, expected interface{}, t *testing.T) {
	t.Errorf("\nactual = %+v\nexpected = %+v", actual, expected)
}

// FailMsg fails and prints actual and expected with a message
func FailMsg(msg string, actual, expected interface{}, t *testing.T) {
	t.Errorf("%s: \nactual = %+v\nexpected = %+v", msg, actual, expected)
}
