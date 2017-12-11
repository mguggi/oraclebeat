package sgastats

import (
	s "github.com/elastic/beats/libbeat/common/schema"
	c "github.com/elastic/beats/libbeat/common/schema/mapstrstr"
)

// Based on https://docs.oracle.com/database/122/REFRN/V-SGASTAT.htm#REFRN30238
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
		"pool":  c.Str("POOL"),
		"name":  c.Str("NAME"),
		"value": c.Str("BYTES"),
	}
)
