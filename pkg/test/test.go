// Contains helpers used for testing
package test

import (
	"testing"
)

// Fail fails and prints actual and expected
func Fail(actual, expected interface{}, t *testing.T) {
	t.Errorf("\nactual = %v\nexpected = %v", actual, expected)
}

// FailMsg fails and prints actual and expected with a message
func FailMsg(msg string, actual, expected interface{}, t *testing.T) {
	t.Errorf(msg, "\nactual = %v\nexpected = %v", actual, expected)
}
