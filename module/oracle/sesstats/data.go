package sesstats

import (
	s "github.com/elastic/beats/libbeat/common/schema"
	c "github.com/elastic/beats/libbeat/common/schema/mapstrstr"
)

// Based on https://docs.oracle.com/database/122/REFRN/V-SESS_TIME_MODEL.htm#REFRN30340
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
		"session": s.Object{
			"id": c.Int("SID"),
		},
		"id":    c.Str("STAT_ID"),
		"name":  c.Str("STAT_NAME"),
		"value": c.Str("VALUE"),
	}
)
