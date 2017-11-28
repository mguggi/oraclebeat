package sesstats

import (
	"testing"

	"github.com/elastic/beats/libbeat/common"
	mbtest "github.com/elastic/beats/metricbeat/mb/testing"
	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	//compose.EnsureUp(t, "oracle")

	f := mbtest.NewEventsFetcher(t, getConfig())
	events, err := f.Fetch()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.True(t, len(events) > 0)
	event := events[0]

	t.Logf("%s/%s event: %+v", f.Module().Name(), f.Name(), event)

	// Check event fields
	database := event["database"].(common.MapStr)
	container := database["container"].(common.MapStr)
	instance := container["instance"].(common.MapStr)
	session := event["session"].(common.MapStr)

	assert.True(t, instance["id"].(int64) > 0)
	assert.True(t, container["id"].(int64) >= 0)
	assert.True(t, session["id"].(int64) > 0)
	assert.NotEmpty(t, event["id"])
	assert.NotEmpty(t, event["name"])
	assert.NotEmpty(t, event["value"])
}

func getConfig() map[string]interface{} {
	return map[string]interface{}{
		"module":     "oracle",
		"metricsets": []string{"sesstats"},
		"hosts":      []string{"oracle://system:oraclebeat@10.1.87.87:1521/ora12"},
	}
}
