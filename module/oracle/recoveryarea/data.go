package recoveryarea

import (
	s "github.com/elastic/beats/libbeat/common/schema"
	c "github.com/elastic/beats/libbeat/common/schema/mapstrstr"
)

// Based on https://docs.oracle.com/en/database/oracle/oracle-database/12.2/refrn/V-RECOVERY_AREA_USAGE.html#GUID-CAF14BAD-9B34-43CD-B297-2A342FA4989B
var (
	schema = s.Schema{
		"database": s.Object{
			"container": s.Object{
				"id": c.Int("CON_ID", s.Optional),
			},
		},
		"file_type": c.Str("FILE_TYPE"),

		"space": s.Object{
			"used": s.Object{
				"pct": c.Float("PERCENT_SPACE_USED"),
			},
			"reclaimable": s.Object{
				"pct": c.Float("PERCENT_SPACE_RECLAIMABLE"),
			},
		},
		"number_of_files": c.Int("NUMBER_OF_FILES"),
	}
)
