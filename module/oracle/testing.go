package oracle

import (
	"os"
)

func GetEnvDSN() string {
	return os.Getenv("ORACLE_DSN")
}
