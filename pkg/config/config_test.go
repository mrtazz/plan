package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigLoading(t *testing.T) {
	assert := assert.New(t)
	tests := map[string]struct {
		filename string
		err      error
	}{
		"default": {
			filename: "./testdata/config.yaml",
			err:      nil,
		},
		"does_not_exist": {
			filename: "./testdata/does-not-exist.yaml",
			err:      ValidationError{message: "open ./testdata/does-not-exist.yaml: no such file or directory"},
		},
		"wrong_date_format": {
			filename: "./testdata/config-wrong-date-format.yaml",
			err:      ValidationError{message: "parsing time \"2023-%m-%d\" as \"2023-%m-%d\": cannot parse \"-%m-%d\" as \"3\""},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := ValidateConfig(tc.filename)
			assert.ErrorIs(tc.err, err)
		})
	}
}
