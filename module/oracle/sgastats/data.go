package sgastats

import (
	s "github.com/elastic/beats/libbeat/common/schema"
	c "github.com/elastic/beats/libbeat/common/schema/mapstrstr"
)

// Based on https://docs.oracle.com/en/database/oracle/oracle-database/12.2/refrn/V-SGASTAT.html#GUID-60D2578E-2293-45F5-91C1-35FDF047E520
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
