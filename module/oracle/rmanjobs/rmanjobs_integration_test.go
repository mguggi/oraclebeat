package rmanjobs

import (
	"testing"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/tests/compose"
	mbtest "github.com/elastic/beats/metricbeat/mb/testing"
	"github.com/mguggi/oraclebeat/module/oracle"
	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	compose.EnsureUpWithTimeout(t, 180, "oracle")

	f := mbtest.NewEventsFetcher(t, getConfig())
	events, err := f.Fetch()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.True(t, len(events) > 0)
	event := events[0]

	// Check event fields
	session := event["session"].(common.MapStr)
	interval := event["interval"].(common.MapStr)

	if session["key"] != nil {
		assert.True(t, session["key"].(int64) > 0)
		assert.True(t, session["recid"].(int64) > 0)
		assert.True(t, session["stamp"].(int64) > 0)
		assert.NotEmpty(t, interval["start"])
		assert.NotEmpty(t, interval["end"])
		assert.NotEmpty(t, interval["duration"])
	}

	if t.Failed() {
		t.Logf("%s/%s event: %+v", f.Module().Name(), f.Name(), event)
	}
}

func getConfig() map[string]interface{} {
	return map[string]interface{}{
		"module":     "oracle",
		"metricsets": []string{"rmanjobs"},
		"hosts":      []string{oracle.GetEnvDSN()},
	}
}
