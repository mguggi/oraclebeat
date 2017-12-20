package rmanjobs

import (
	"time"

	s "github.com/elastic/beats/libbeat/common/schema"
	c "github.com/elastic/beats/libbeat/common/schema/mapstrstr"
)

// Based on https://docs.oracle.com/en/database/oracle/oracle-database/12.2/refrn/V-RMAN_BACKUP_JOB_DETAILS.html#GUID-6DD6941F-C612-4EB1-BE34-0CB5ECF331A5
var (
	schema = s.Schema{
		"database": s.Object{
			"container": s.Object{
				"id": c.Int("CON_ID", s.Optional),
			},
		},
		"session": s.Object{
			"key":   c.Int("SESSION_KEY"),
			"recid": c.Int("SESSION_RECID"),
			"stamp": c.Int("SESSION_STAMP"),
		},
		"command_id": c.Str("COMMAND_ID"),
		"interval": s.Object{
			"start":    c.Time(time.RFC3339, "START_TIME"),
			"end":      c.Time(time.RFC3339, "END_TIME"),
			"duration": c.Float("ELAPSED_SECONDS"),
		},
		"input": s.Object{
			"type":          c.Str("INPUT_TYPE"),
			"bytes":         c.Int("INPUT_BYTES"),
			"bytes_per_sec": c.Float("INPUT_BYTES_PER_SEC"),
		},
		"output": s.Object{
			"device_type":   c.Str("OUTPUT_DEVICE_TYPE"),
			"bytes":         c.Int("OUTPUT_BYTES"),
			"bytes_per_sec": c.Float("OUTPUT_BYTES_PER_SEC"),
		},
		"backed_by_osb": c.Str("BACKED_BY_OSB"),
		"autobackup": s.Object{
			"count": c.Int("AUTOBACKUP_COUNT"),
			"done":  c.Str("AUTOBACKUP_DONE"),
		},
		"status":            c.Str("STATUS"),
		"optimized":         c.Str("OPTIMIZED"),
		"compression_ratio": c.Float("COMPRESSION_RATIO"),
	}
)
