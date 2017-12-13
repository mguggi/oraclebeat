package sqlstats

import (
	"time"

	s "github.com/elastic/beats/libbeat/common/schema"
	c "github.com/elastic/beats/libbeat/common/schema/mapstrstr"
)

// Based on https://docs.oracle.com/en/database/oracle/oracle-database/12.2/refrn/V-SQLSTATS.html#GUID-495DD17D-6741-433F-871D-C965EB221DA9
var (
	schema = s.Schema{
		"database": s.Object{
			"container": s.Object{
				"id":   c.Int("CON_ID", s.Optional),
				"dbid": c.Int("CON_DBID", s.Optional),
				"instance": s.Object{
					"id": c.Int("INST_ID"),
				},
			},
		},
		"sql": s.Object{
			"id":   c.Str("SQL_ID"),
			"text": c.Str("SQL_FULLTEXT"),
		},
		"last_active": s.Object{
			"time":          c.Time(time.RFC3339, "LAST_ACTIVE_TIME", s.Optional),
			"child_address": c.Str("LAST_ACTIVE_CHILD_ADDRESS"),
		},
		"plan_hash_value": c.Int("PLAN_HASH_VALUE"),
		"parse_calls":     c.Int("PARSE_CALLS"),
		"disk_reads":      c.Int("DISK_READS"),
		"direct": s.Object{
			"writes": c.Int("DIRECT_WRITES"),
			"reads":  c.Int("DIRECT_READS", s.Optional),
		},
		"buffer_gets":           c.Int("BUFFER_GETS"),
		"rows_processed":        c.Int("ROWS_PROCESSED"),
		"serializable_aborts":   c.Int("SERIALIZABLE_ABORTS"),
		"fetches":               c.Int("FETCHES"),
		"executions":            c.Int("EXECUTIONS"),
		"end_of_fetch_count":    c.Int("END_OF_FETCH_COUNT"),
		"loads":                 c.Int("LOADS"),
		"version_count":         c.Int("VERSION_COUNT"),
		"invalidations":         c.Int("INVALIDATIONS"),
		"px_servers_executions": c.Int("PX_SERVERS_EXECUTIONS"),
		"cpu_time":              c.Int("CPU_TIME"),
		"elapsed_time":          c.Int("ELAPSED_TIME"),
		"avg_hard_parse_time":   c.Int("AVG_HARD_PARSE_TIME"),
		"application_wait_time": c.Int("APPLICATION_WAIT_TIME"),
		"concurrency_wait_time": c.Int("CONCURRENCY_WAIT_TIME"),
		"cluster_wait_time":     c.Int("CLUSTER_WAIT_TIME"),
		"user_io_wait_time":     c.Int("USER_IO_WAIT_TIME"),
		"plsql_exec_time":       c.Int("PLSQL_EXEC_TIME"),
		"java_exec_time":        c.Int("JAVA_EXEC_TIME"),
		"sorts":                 c.Int("SORTS"),
		"sharable_mem":          c.Int("SHARABLE_MEM"),
		"total_sharable_mem":    c.Int("TOTAL_SHARABLE_MEM"),
		"typecheck_mem":         c.Int("TYPECHECK_MEM"),
		"io_interconnect_bytes": c.Int("IO_INTERCONNECT_BYTES"),
		"physical": s.Object{
			"read": s.Object{
				"requests": c.Int("PHYSICAL_READ_REQUESTS"),
				"bytes":    c.Int("PHYSICAL_READ_BYTES"),
			},
			"write": s.Object{
				"requests": c.Int("PHYSICAL_WRITE_REQUESTS"),
				"bytes":    c.Int("PHYSICAL_WRITE_BYTES"),
			},
		},
		"exact_matching_signature": c.Str("EXACT_MATCHING_SIGNATURE"),
		"force_matching_signature": c.Str("FORCE_MATCHING_SIGNATURE"),
		"io_cell": s.Object{
			"offload_eligible_bytes": c.Int("IO_CELL_OFFLOAD_ELIGIBLE_BYTES"),
			"uncompressed_bytes":     c.Int("IO_CELL_UNCOMPRESSED_BYTES"),
			"offload_returned_bytes": c.Int("IO_CELL_OFFLOAD_RETURNED_BYTES"),
		},
		"obsolete_count": c.Int("OBSOLETE_COUNT", s.Optional),
	}
)
