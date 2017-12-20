package oracle

import (
	"database/sql"
	"fmt"
	"hash/fnv"

	"encoding/hex"

	"bytes"
	"io"

	"time"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/metricbeat/mb"
	"github.com/pkg/errors"
	"gopkg.in/goracle.v2"
)

const selector = "oracle"

func init() {
	// Register the ModuleFactory function for the "oracle" module.
	if err := mb.Registry.AddModule("oracle", NewModule); err != nil {
		panic(err)
	}
}

func NewModule(base mb.BaseModule) (mb.Module, error) {
	// Validate that at least one host has been specified.
	config := struct {
		Hosts []string `config:"hosts"     validate: "nonzero,required"`
	}{}

	if err := base.UnpackConfig(&config); err != nil {
		return nil, err
	}
	return &base, nil
}

// ParseDSN creates a DSN (data source name) string by parsing the host.
// It validates the resulting DSN and returns an error if the DSN is invalid.
//
//   Format:  username:password@address:port/service
//   Example: oracle://system:changeme@hostname:1521/orcl
func ParseDSN(mod mb.Module, host string) (mb.HostData, error) {

	// Parse the provided connection string.
	cp, err := goracle.ParseConnString(host)
	if err != nil {
		return mb.HostData{}, errors.Wrapf(err, "%s failed parsing url", selector)
	}
	//username, password, service := ora.SplitDSN(host)
	logp.Debug(selector, "%s:***@%s", cp.Username, cp.SID)

	return mb.HostData{
		URI:          host,
		SanitizedURI: fmt.Sprintf("%s:***@%s", cp.Username, cp.SID),
		Host:         cp.SID,
		User:         cp.Username,
		Password:     cp.Password,
	}, nil
}

// Scan data retrived from metricset queries
func Scan(rows *sql.Rows) ([]map[string]interface{}, error) {

	columns, err := rows.ColumnTypes()
	if err != nil {
		return nil, errors.Wrapf(err, "%s failed scanning columns", selector)
	}

	vals := make([]interface{}, len(columns))
	valsPtr := make([]interface{}, len(columns))
	for i := range columns {
		valsPtr[i] = &vals[i]
	}

	results := []map[string]interface{}{}
	for rows.Next() {
		err = rows.Scan(valsPtr...)
		if err != nil {
			return nil, errors.Wrapf(err, "%s failed scanning row", selector)
		}

		var buf bytes.Buffer
		result := map[string]interface{}{}
		for i, _ := range columns {
			switch columns[i].DatabaseTypeName() {
			case "CLOB":
				lob := vals[i].(*goracle.Lob)
				buf.Reset()
				if _, err := io.Copy(&buf, lob); err != nil {
					return nil, err
				}
				result[columns[i].Name()] = buf.String()
			case "NUMBER":
				if vals[i] != nil {
					result[columns[i].Name()] = vals[i].(goracle.Number).String()
				}
			case "DATE":
				result[columns[i].Name()] = vals[i].(time.Time).Format(time.RFC3339)
			case "TIMESTAMP WITH TIMEZONE":
				result[columns[i].Name()] = vals[i].(time.Time).Format(time.RFC3339Nano)
			case "RAW":
				result[columns[i].Name()] = hex.EncodeToString(vals[i].([]byte))
			default:
				result[columns[i].Name()] = vals[i].(string)
			}
		}

		logp.Debug(selector, "result: %v", result)
		results = append(results, result)
	}
	return results, nil
}

// Generates a short unique hash and returns the result as hex string
func GetHash(text string) string {
	hasher := fnv.New32()
	hasher.Write([]byte(text))
	return fmt.Sprintf("%x", hasher.Sum32())
}
