package db

import (
	"github.com/eugeneverywhere/billing/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCases = map[string]struct {
	input  config.DBConfig
	result string
}{
	"standart generating db url": {
		config.DBConfig{
			Host:     "localhost",
			Port:     3306,
			User:     "root",
			Password: "111",
			Name:     "test",
		},
		"root:111@tcp(localhost:3306)/test?multiStatements=true&parseTime=true",
	},
}

func TestGenerateMySQLDatabaseURL(t *testing.T) {
	for name, test := range testCases {
		t.Logf("DB tests: Running test case '%s'", name)

		res := GenerateMySQLDatabaseURL(test.input)
		assert.Equal(t, test.result, res)
	}
}
