package filespace

import (
	"context"
	"database/sql"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/metricbeat/mb"
	"github.com/mguggi/oraclebeat/module/oracle"
	"github.com/pkg/errors"

	_ "gopkg.in/goracle.v2"
)

const selector = "oracle/filespace"

// init registers the MetricSet with the central registry.
// The New method will be called after the setup of the module and before starting to fetch data
func init() {
	if err := mb.Registry.AddMetricSet("oracle", "filespace", New, oracle.ParseDSN); err != nil {
		panic(err)
	}
}

// MetricSet type defines all fields of the MetricSet
// As a minimum it must inherit the mb.BaseMetricSet fields, but can be extended with
// additional entries. These variables can be used to persist data or configuration between
// multiple fetch calls.
type MetricSet struct {
	mb.BaseMetricSet
	timeout time.Duration
	db      *sql.DB
	stmt    *sql.Stmt
}

// New create a new instance of the MetricSet
// Part of new is also setting up the configuration by processing additional
// configuration entries if needed.
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {

	config := struct{}{}

	logp.Info("EXPERIMENTAL: The oracle filespace metricset is experimental")

	if err := base.Module().UnpackConfig(&config); err != nil {
		return nil, err
	}

	// Configure the driver to return large floats as strings
	logp.Debug(selector, "Connect to host: %s", base.HostData().URI)
	db, err := sql.Open("goracle", base.HostData().URI)
	if err != nil {
		return nil, errors.Wrapf(err, "%s open connection failed", selector)
	}

	stmt, err := db.Prepare(
		"select /* metricset: filespace */ " +
			"ts.name, decode(ts.bigfile, 'NO', 'F', 'T') bigfile, decode(ts.flashback_on, 'NO', 'F', 'T') flashback_on, fsu.* " +
			"from v$tablespace ts " +
			"join v$filespace_usage fsu on ts.ts# = fsu.tablespace_id")

	if err != nil {
		return nil, errors.Wrapf(err, "%s failed prepare statement", selector)
	}

	return &MetricSet{
		BaseMetricSet: base,
		db:            db,
		stmt:          stmt,
	}, nil
}

// Fetch methods implements the data gathering and data conversion to the right format
// It returns the event which is then forward to the output. In case of an error, a
// descriptive error must be returned.
func (m *MetricSet) Fetch() ([]common.MapStr, error) {

	// Make connection timeout aware
	timeout := m.BaseMetricSet.Module().Config().Timeout
	logp.Debug(selector, "Set timout to: %s", timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	rows, err := m.stmt.QueryContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "%s failed query data", selector)
	}
	defer rows.Close()

	events := []common.MapStr{}
	results, _ := oracle.Scan(rows)
	for _, result := range results {
		data, _ := schema.Apply(result)
		events = append(events, data)
	}

	return events, nil
}

// Closer is an optional interface that a MetricSet can implement in order to
// cleanup any resources it has open at shutdown.
func (m *MetricSet) Close() error {
	defer m.db.Close()
	defer m.stmt.Close()
	return nil
}
