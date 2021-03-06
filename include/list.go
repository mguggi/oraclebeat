/*
Package include imports all Module and MetricSet packages so that they register
their factories with the global registry. This package can be imported in the
main package to automatically register all of the standard supported Metricbeat
modules.
*/
package include

import (
	// This list is automatically generated by `make imports`
	_ "github.com/mguggi/oraclebeat/module/oracle"
	_ "github.com/mguggi/oraclebeat/module/oracle/alertlog"
	_ "github.com/mguggi/oraclebeat/module/oracle/osstats"
	_ "github.com/mguggi/oraclebeat/module/oracle/recoveryarea"
	_ "github.com/mguggi/oraclebeat/module/oracle/rmanjobs"
	_ "github.com/mguggi/oraclebeat/module/oracle/session"
	_ "github.com/mguggi/oraclebeat/module/oracle/sessmetric"
	_ "github.com/mguggi/oraclebeat/module/oracle/sesstats"
	_ "github.com/mguggi/oraclebeat/module/oracle/sgastats"
	_ "github.com/mguggi/oraclebeat/module/oracle/sqlstats"
	_ "github.com/mguggi/oraclebeat/module/oracle/sysmetric"
	_ "github.com/mguggi/oraclebeat/module/oracle/sysstats"
	_ "github.com/mguggi/oraclebeat/module/oracle/tablespace"
)
