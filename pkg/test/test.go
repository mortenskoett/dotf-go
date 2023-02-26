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

// Fail fails and prints actual and expected
func Fail(actual, expected interface{}, t *testing.T) {
	t.Errorf("\nactual = %v\nexpected = %v", actual, expected)
}

// FailMsg fails and prints actual and expected with a message
func FailMsg(msg string, actual, expected interface{}, t *testing.T) {
	t.Errorf("%s: \nactual = %v\nexpected = %v", msg, actual, expected)
}
