package sysstats

import (
	s "github.com/elastic/beats/libbeat/common/schema"
	c "github.com/elastic/beats/libbeat/common/schema/mapstrstr"
)

// Based on https://docs.oracle.com/en/database/oracle/oracle-database/12.2/refrn/V-SYS_TIME_MODEL.html#GUID-DC16AB84-4978-497B-8AFB-C3C23D83FC3C
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
