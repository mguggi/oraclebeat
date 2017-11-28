package filespace

import (
	s "github.com/elastic/beats/libbeat/common/schema"
	c "github.com/elastic/beats/libbeat/common/schema/mapstrstr"
)

// Based on https://docs.oracle.com/database/122/REFRN/V-TABLESPACE.htm#REFRN30277
//          https://docs.oracle.com/database/122/REFRN/V-FILESPACE_USAGE.htm#REFRN30333
var (
	schema = s.Schema{
		"database": s.Object{
			"container": s.Object{
				"id": c.Int("CON_ID"),
			},
		},
		"tablespace": s.Object{
			"id":        c.Int("TABLESPACE_ID"),
			"name":      c.Str("NAME"),
			"bigfile":   c.Bool("BIGFILE"),
			"flashback": c.Bool("FLASHBACK_ON"),
		},
		"rfno":            c.Int("RFNO"),
		"allocated_space": c.Int("ALLOCATED_SPACE"),
		"file_size":       c.Int("FILE_SIZE"),
		"file_maxsize":    c.Int("FILE_SIZE"),
		"changescn": s.Object{
			"base": c.Str("CHANGESCN_BASE"),
			"wrap": c.Str("CHANGESCN_WRAP"),
		},
		"flag": c.Int("FLAG"),
	}
)
