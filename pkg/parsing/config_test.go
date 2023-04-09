package parsing_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mortenskoett/dotf-go/pkg/parsing"
)

func Test_NewConfigMap_returns_valid_map_configuration(t *testing.T) {
	conf := parsing.NewSensibleConfiguration()
	res, err := parsing.ConvertConfigToMap(conf)
	if err != nil {
		t.Error(res, err)
	}
}

func Test_CreateSerializableConfig_returns_valid_config_as_bytes(t *testing.T) {
	testinput := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	expected := fmt.Sprintf("key1 = value1\nkey2 = value2\n")

	bs := parsing.CreateSerializableConfig(testinput)
	s := string(bs)

	diff := cmp.Diff(s, expected)
	if diff != "" {
		t.Errorf("have: %+v\nwant: %+v\ndiff: %+v", s, expected, diff)
	}
}
