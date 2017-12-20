package alertlog

import (
	"time"

	s "github.com/elastic/beats/libbeat/common/schema"
	c "github.com/elastic/beats/libbeat/common/schema/mapstrstr"
)

// Based on https://docs.oracle.com/en/database/oracle/oracle-database/12.2/refrn/V-DIAG_ALERT_EXT.html#GUID-7EC93FE0-FF30-4A94-92BC-785E2BCB38F3
var (
	schema = s.Schema{
		"database": s.Object{
			"container": s.Object{
				"id":   c.Int("CON_ID", s.Optional),
				"uid":  c.Int("CON_UID", s.Optional),
				"name": c.Str("CONTAINER_NAME", s.Optional),
				"instance": s.Object{
					"id": c.Int("INST_ID"),
				},
			},
		},
		"timestamp":    c.Time(time.RFC3339Nano, "ORIGINATING_TIMESTAMP"),
		"adr_home":     c.Str("ADR_HOME"),
		"organization": c.Str("ORGANIZATION_ID"),
		"component":    c.Str("COMPONENT_ID"),
		"host": s.Object{
			"id":      c.Str("HOST_ID"),
			"address": c.Str("HOST_ADDRESS"),
		},
		"message": s.Object{
			"id":         c.Str("MESSAGE_ID", s.Optional),
			"type":       c.Int("MESSAGE_TYPE"),
			"level":      c.Int("MESSAGE_LEVEL"),
			"group":      c.Str("MESSAGE_GROUP", s.Optional),
			"text":       c.Str("MESSAGE_TEXT"),
			"arguments":  c.Str("MESSAGE_ARGUMENTS", s.Optional),
			"attributes": c.Str("SUPPLEMENTAL_ATTRIBUTES", s.Optional),
			"details":    c.Str("SUPPLEMENTAL_DETAILS", s.Optional),
			"record":     c.Int("RECORD_ID"),
			"version":    c.Str("VERSION"),
		},
		"client":   c.Str("CLIENT_ID", s.Optional),
		"module":   c.Str("MODULE_ID", s.Optional),
		"process":  c.Str("PROCESS_ID", s.Optional),
		"thread":   c.Str("THREAD_ID", s.Optional),
		"user":     c.Str("USER_ID", s.Optional),
		"partion":  c.Int("PARTITION"),
		"filename": c.Str("FILENAME"),
		"problem":  c.Str("PROBLEM_KEY", s.Optional),
	}
)
