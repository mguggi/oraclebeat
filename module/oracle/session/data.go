package session

import (
	"time"

	s "github.com/elastic/beats/libbeat/common/schema"
	c "github.com/elastic/beats/libbeat/common/schema/mapstrstr"
)

// Based on https://docs.oracle.com/en/database/oracle/oracle-database/12.2/refrn/V-SESSION.html#GUID-28E2DC75-E157-4C0A-94AB-117C205789B9
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
		"address":          c.Str("SADDR"),
		"id":               c.Int("SID"),
		"serial_number":    c.Int("SERIAL#"),
		"audit_session_id": c.Int("AUDSID", s.Optional),
		"status":           c.Str("STATUS"),
		"type":             c.Str("TYPE"),
		"process": s.Object{
			"id":       c.Str("PROCESS"),
			"addresss": c.Str("PADDR"),
			"user":     c.Str("OSUSER"),
			"machine":  c.Str("MACHINE"),
			"port":     c.Int("PORT"),
			"terminal": c.Str("TERMINAL"),
			"program":  c.Str("PROGRAM"),
			"type":     c.Str("SERVER"),
		},
		"user": s.Object{
			"id":   c.Int("USER#", s.Optional),
			"name": c.Str("USERNAME"),
			"schema": s.Object{
				"id":   c.Int("SCHEMA#"),
				"name": c.Str("SCHEMANAME"),
			},
		},
		"command": c.Int("COMMAND"),
		"owner":   c.Int("OWNERID", s.Optional),
		"transaction": s.Object{
			"address": c.Str("TADDR", s.Optional),
			"lock": s.Object{
				"address": c.Str("LOCKWAIT", s.Optional),
			},
		},
		"sql": s.Object{
			"address":      c.Str("SQL_ADDRESS"),
			"hash":         c.Int("SQL_HASH_VALUE", s.Optional),
			"id":           c.Str("SQL_ID", s.Optional),
			"child_number": c.Int("SQL_CHILD_NUMBER", s.Optional),
			"execution": s.Object{
				"start": c.Time(time.RFC3339, "SQL_EXEC_START", s.Optional),
				"id":    c.Int("SQL_EXEC_ID", s.Optional),
			},
			"precursor": s.Object{
				"address":      c.Str("PREV_SQL_ADDR"),
				"hash":         c.Int("PREV_HASH_VALUE", s.Optional),
				"id":           c.Str("PREV_SQL_ID", s.Optional),
				"child_number": c.Int("PREV_CHILD_NUMBER", s.Optional),
				"execution": s.Object{
					"start": c.Time(time.RFC3339, "PREV_EXEC_START", s.Optional),
					"id":    c.Int("PREV_EXEC_ID", s.Optional),
				},
			},
		},
		"plsql": s.Object{
			"entry": s.Object{
				"object_id":     c.Int("PLSQL_ENTRY_OBJECT_ID", s.Optional),
				"subprogram_id": c.Int("PLSQL_ENTRY_SUBPROGRAM_ID", s.Optional),
			},
			"object_id":          c.Int("PLSQL_OBJECT_ID", s.Optional),
			"subprogram_id":      c.Int("PLSQL_SUBPROGRAM_ID", s.Optional),
			"debugger_connected": c.Bool("PLSQL_DEBUGGER_CONNECTED", s.Optional),
		},
		"module": s.Object{
			"name": c.Str("MODULE", s.Optional),
			"hash": c.Str("MODULE_HASH", s.Optional),
		},
		"action": s.Object{
			"name": c.Str("ACTION", s.Optional),
			"hash": c.Str("ACTION_HASH", s.Optional),
		},
		"client": s.Object{
			"id":   c.Str("CLIENT_IDENTIFIER"),
			"info": c.Str("CLIENT_INFO"),
		},
		"fixed_table_sequence": c.Int("FIXED_TABLE_SEQUENCE", s.Optional),
		"row_lock": s.Object{
			"object_id": c.Int("ROW_WAIT_OBJ#"),
			"file_id":   c.Int("ROW_WAIT_FILE#"),
			"block_id":  c.Int("ROW_WAIT_BLOCK#"),
			"row_id":    c.Int("ROW_WAIT_ROW#"),
		},
		"top_level_call_number": c.Int("TOP_LEVEL_CALL_NUMBER", s.Optional),
		"logon_time":            c.Time(time.RFC3339, "LOGON_TIME"),
		"last_call_et":          c.Int("LAST_CALL_ET", s.Optional),
		"failover": s.Object{
			"type":        c.Str("FAILOVER_TYPE"),
			"method":      c.Str("FAILOVER_METHOD"),
			"failed_over": c.Str("FAILED_OVER"),
		},
		"resource_consumer_group": c.Str("RESOURCE_CONSUMER_GROUP"),
		"pdml_status":             c.Str("PDML_STATUS"),
		"pddl_status":             c.Str("PDDL_STATUS"),
		"pq_status":               c.Str("PQ_STATUS"),
		"current_queue_duration":  c.Int("CURRENT_QUEUE_DURATION"),
		"blocker": s.Object{
			"status":      c.Str("BLOCKING_SESSION_STATUS"),
			"instance_id": c.Int("BLOCKING_INSTANCE", s.Optional),
			"session_id":  c.Int("BLOCKING_SESSION", s.Optional),
			"final": s.Object{
				"status":      c.Str("FINAL_BLOCKING_SESSION_STATUS"),
				"instance_id": c.Int("FINAL_BLOCKING_INSTANCE", s.Optional),
				"session_id":  c.Int("FINAL_BLOCKING_SESSION", s.Optional),
			},
		},
		"wait": s.Object{
			"sequence": c.Int("SEQ#"),
			"event": s.Object{
				"number": c.Int("EVENT#"),
				"value":  c.Str("EVENT"),
			},
			"p1": s.Object{
				"number": c.Int("P1", s.Optional),
				"text":   c.Str("P1TEXT"),
				"raw":    c.Str("P1RAW"),
			},
			"p2": s.Object{
				"number": c.Int("P2", s.Optional),
				"text":   c.Str("P2TEXT"),
				"raw":    c.Str("P2RAW"),
			},
			"p3": s.Object{
				"number": c.Int("P3", s.Optional),
				"text":   c.Str("P3TEXT"),
				"raw":    c.Str("P3RAW"),
			},
			"class": s.Object{
				"id":     c.Int("WAIT_CLASS_ID", s.Optional),
				"number": c.Int("WAIT_CLASS#"),
				"value":  c.Str("WAIT_CLASS"),
			},
			"state":                   c.Str("STATE"),
			"time_us":                 c.Int("WAIT_TIME_MICRO", s.Optional),
			"time_remaining_us":       c.Int("TIME_REMAINING_MICRO", s.Optional),
			"time_since_last_wait_us": c.Int("TIME_SINCE_LAST_WAIT_MICRO", s.Optional),
		},
		"service_name": c.Str("SERVICE_NAME"),
		"trace": s.Object{
			"sql":        c.Str("SQL_TRACE"),
			"waits":      c.Bool("SQL_TRACE_WAITS"),
			"binds":      c.Bool("SQL_TRACE_BINDS"),
			"plan_stats": c.Str("SQL_TRACE_PLAN_STATS"),
		},
		"edition_id": c.Int("SESSION_EDITION_ID"),
		"creator": s.Object{
			"address":       c.Str("CREATOR_ADDR"),
			"serial_number": c.Int("CREATOR_SERIAL#"),
		},
		"execution_context_id":       c.Str("ECID"),
		"sql_translation_profile_id": c.Int("SQL_TRANSLATION_PROFILE_ID", s.Optional),
		"pga_tunable_mem":            c.Int("PGA_TUNABLE_MEM", s.Optional),
		"shard_ddl_status":           c.Str("SHARD_DDL_STATUS", s.Optional),
		"external_name":              c.Str("EXTERNAL_NAME", s.Optional),
	}
)
