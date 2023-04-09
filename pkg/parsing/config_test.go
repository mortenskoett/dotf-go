package parsing_test

import (
	"testing"

	"github.com/mortenskoett/dotf-go/pkg/parsing"
)

func Test_NewConfigMap_returns_valid_map_configuration(t *testing.T) {
	res, err := parsing.NewConfigMap()
	if err != nil {
		t.Error(res, err)
	}
}
