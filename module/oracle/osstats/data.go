package osstats

import (
	s "github.com/elastic/beats/libbeat/common/schema"
	c "github.com/elastic/beats/libbeat/common/schema/mapstrstr"
)

// Based on https://docs.oracle.com/en/database/oracle/oracle-database/12.2/refrn/V-OSSTAT.html#GUID-E1E48692-47FA-4AE3-9402-82477E66FFC0
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
		"id":         c.Int("OSSTAT_ID"),
		"name":       c.Str("STAT_NAME"),
		"value":      c.Str("VALUE"),
		"comments":   c.Str("COMMENTS"),
		"cumulative": c.Str("CUMULATIVE"),
	}
)
