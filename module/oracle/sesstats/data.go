package sesstats

import (
	s "github.com/elastic/beats/libbeat/common/schema"
	c "github.com/elastic/beats/libbeat/common/schema/mapstrstr"
)

// https://docs.oracle.com/en/database/oracle/oracle-database/12.2/refrn/V-SESS_TIME_MODEL.html#GUID-B5CF4362-325D-4F22-9A08-0873FA32A5C0
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
