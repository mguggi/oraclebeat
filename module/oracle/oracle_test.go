package oracle

import (
	"testing"

	mbtest "github.com/elastic/beats/metricbeat/mb/testing"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDSN(t *testing.T) {
	tests := []struct {
		Name     string
		Host     string
		Timeout  time.Duration
		Expected string
	}{
		{
			Name:     "complete url",
			Host:     "oracle://user:passwd@host:1521/orcl",
			Expected: "host:1521/orcl",
		},
		{
			Name:     "url without port",
			Host:     "oracle://user:passwd@host/orcl",
			Expected: "host/orcl",
		},
	}

	for _, test := range tests {
		mod := mbtest.NewTestModule(t, map[string]interface{}{})
		mod.ModConfig.Timeout = test.Timeout
		hostData, err := ParseDSN(mod, test.Host)
		if err != nil {
			t.Error(err)
			continue
		}

		assert.Equal(t, test.Expected, hostData.Host, test.Name)
	}
}
