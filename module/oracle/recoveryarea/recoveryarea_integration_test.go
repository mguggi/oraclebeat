package recoveryarea

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
	space := event["space"].(common.MapStr)
	used := space["used"].(common.MapStr)
	reclaimable := space["reclaimable"].(common.MapStr)

	assert.NotEmpty(t, event["file_type"])
	assert.True(t, used["pct"].(float64) >= 0)
	assert.True(t, reclaimable["pct"].(float64) >= 0)
	assert.True(t, event["number_of_files"].(int64) >= 0)

	if t.Failed() {
		t.Logf("%s/%s event: %+v", f.Module().Name(), f.Name(), event)
	}
}

func getConfig() map[string]interface{} {
	return map[string]interface{}{
		"module":     "oracle",
		"metricsets": []string{"recoveryarea"},
		"hosts":      []string{oracle.GetEnvDSN()},
	}
}
