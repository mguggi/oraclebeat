package sysmetric

import (
	"time"

	s "github.com/elastic/beats/libbeat/common/schema"
	c "github.com/elastic/beats/libbeat/common/schema/mapstrstr"
)

// Based on https://docs.oracle.com/database/122/REFRN/V-SYSMETRIC.htm#REFRN30343
var (
	schema = s.Schema{
		"database": s.Object{
			"container": s.Object{
				"id": c.Int("CON_ID", s.Optional),
				"instance": s.Object{
					"id": c.Int("INST_ID"),
				},
			},
		},
		"interval": s.Object{
			"begin": c.Time(time.RFC3339, "BEGIN_TIME"),
			"end":   c.Time(time.RFC3339, "END_TIME"),
			"csec":  c.Int("INTSIZE_CSEC"),
		},
		"group_id": c.Int("GROUP_ID"),
		"id":       c.Int("METRIC_ID"),
		"name":     c.Str("METRIC_NAME"),
		"value":    c.Float("VALUE"),
		"unit":     c.Str("METRIC_UNIT"),
	}
)
