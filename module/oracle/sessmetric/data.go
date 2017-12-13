package sessmetric

import (
	"time"

	s "github.com/elastic/beats/libbeat/common/schema"
	c "github.com/elastic/beats/libbeat/common/schema/mapstrstr"
)

// Based on https://docs.oracle.com/en/database/oracle/oracle-database/12.2/refrn/V-SESSMETRIC.html#GUID-FFC5A6B9-9FBE-477A-91A8-0675C345725C
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
			"id":            c.Int("SESSION_ID"),
			"serial_number": c.Int("SERIAL_NUM"),
		},
		"interval": s.Object{
			"begin": c.Time(time.RFC3339, "BEGIN_TIME"),
			"end":   c.Time(time.RFC3339, "END_TIME"),
			"csec":  c.Int("INTSIZE_CSEC"),
		},
		"cpu": c.Str("CPU"),
		"reads": s.Object{
			"physical": s.Object{
				"count": c.Int("PHYSICAL_READS"),
				"pct":   c.Float("PHYSICAL_READ_PCT"),
			},
			"logical": s.Object{
				"count": c.Int("LOGICAL_READS"),
				"pct":   c.Float("LOGICAL_READ_PCT"),
			},
		},
		"pga_memory": c.Str("PGA_MEMORY"),
		"parses": s.Object{
			"hard": c.Int("HARD_PARSES"),
			"soft": c.Int("SOFT_PARSES"),
		},
	}
)
