package alertlog

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
	message := event["message"].(common.MapStr)
	host := event["host"].(common.MapStr)

	assert.NotEmpty(t, event["timestamp"])
	assert.NotEmpty(t, event["filename"])
	assert.NotEmpty(t, event["adr_home"])
	assert.NotEmpty(t, event["organization"])
	assert.NotEmpty(t, event["component"])
	assert.NotEmpty(t, host["id"])
	assert.NotEmpty(t, host["address"])
	assert.NotEmpty(t, message["type"])
	assert.NotEmpty(t, message["level"])
	assert.NotEmpty(t, message["text"])
	assert.True(t, message["record"].(int64) > 0)

	if !t.Failed() {
		t.Logf("%s/%s event: %+v", f.Module().Name(), f.Name(), event)
	}
}

func getConfig() map[string]interface{} {
	return map[string]interface{}{
		"module":     "oracle",
		"metricsets": []string{"alertlog"},
		"hosts":      []string{oracle.GetEnvDSN()},
	}
}
