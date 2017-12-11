package tablespace

import (
	s "github.com/elastic/beats/libbeat/common/schema"
	c "github.com/elastic/beats/libbeat/common/schema/mapstrstr"
)

// Based on https://docs.oracle.com/database/122/REFRN/DBA_TABLESPACE_USAGE_METRICS.htm#REFRN23496
var (
	schema = s.Schema{
		"name": c.Str("TABLESPACE_NAME"),

		"size": c.Int("TABLESPACE_SIZE"),
		"used": s.Object{
			"space": c.Int("USED_SPACE"),
			"pct":   c.Float("USED_PERCENT"),
		},
	}
)
