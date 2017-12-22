package sessmetric

import (
	"testing"

	"github.com/mguggi/oraclebeat/module/oracle"
	"github.com/stretchr/testify/assert"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/tests/compose"
	mbtest "github.com/elastic/beats/metricbeat/mb/testing"
)

func TestFetch(t *testing.T) {
	compose.EnsureUpWithTimeout(t, 180, "oracle")

	f := mbtest.NewEventsFetcher(t, getConfig())
	events, err := f.Fetch()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	if len(events) > 0 {
		event := events[0]

		// Check event fields
		database := event["database"].(common.MapStr)
		container := database["container"].(common.MapStr)
		instance := container["instance"].(common.MapStr)
		session := event["session"].(common.MapStr)
		interval := event["interval"].(common.MapStr)

		assert.True(t, instance["id"].(int64) > 0)
		assert.True(t, session["id"].(int64) > 0)
		assert.True(t, session["serial_number"].(int64) > 0)
		assert.NotEmpty(t, interval["begin"])
		assert.NotEmpty(t, interval["end"])
		assert.NotEmpty(t, event["cpu"])
		assert.NotEmpty(t, event["pga_memory"])

		if t.Failed() {
			t.Logf("%s/%s event: %+v", f.Module().Name(), f.Name(), event)
		}
	}
}

func getConfig() map[string]interface{} {
	return map[string]interface{}{
		"module":     "oracle",
		"metricsets": []string{"sessmetric"},
		"hosts":      []string{oracle.GetEnvDSN()},
	}
}
