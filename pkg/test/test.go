// Contains helpers used for testing
package test

import (
	"testing"
)

// Helper that fails and prints actual and expected
func Fail(actual, expected interface{}, t *testing.T) {
	t.Errorf("\nactual = %v\nexpected = %v", actual, expected)
}
