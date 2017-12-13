package tablespace

import (
	s "github.com/elastic/beats/libbeat/common/schema"
	c "github.com/elastic/beats/libbeat/common/schema/mapstrstr"
)

// Based on https://docs.oracle.com/en/database/oracle/oracle-database/12.2/refrn/DBA_TABLESPACE_USAGE_METRICS.html#GUID-FE479528-BB37-4B55-92CF-9EC19EDF4F46
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
