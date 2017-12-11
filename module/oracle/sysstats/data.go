package sysstats

import (
	s "github.com/elastic/beats/libbeat/common/schema"
	c "github.com/elastic/beats/libbeat/common/schema/mapstrstr"
)

// Based on https://docs.oracle.com/database/122/REFRN/V-SYS_TIME_MODEL.htm#REFRN30346
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
		"id":    c.Str("STAT_ID"),
		"name":  c.Str("STAT_NAME"),
		"value": c.Str("VALUE"),
	}
)
