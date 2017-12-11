package sysstats

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
	database := event["database"].(common.MapStr)
	container := database["container"].(common.MapStr)
	instance := container["instance"].(common.MapStr)

	assert.True(t, instance["id"].(int64) > 0)
	assert.NotEmpty(t, event["id"])
	assert.NotEmpty(t, event["name"])
	assert.NotEmpty(t, event["value"])

	if t.Failed() {
		t.Logf("%s/%s event: %+v", f.Module().Name(), f.Name(), event)
	}
}

func getConfig() map[string]interface{} {
	return map[string]interface{}{
		"module":     "oracle",
		"metricsets": []string{"sysstats"},
		"hosts":      []string{oracle.GetEnvDSN()},
	}
}
