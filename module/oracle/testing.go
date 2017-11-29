package oracle

import (
	"os"
)

func GetEnvDSN() string {
        dsn := os.Getenv("ORACLE_DSN")

	if len(dsn) == 0 {
		dsn = "oracle://system:chang3Me@localhost:1521/ORCLPDB1"
	}

	return dsn
}
