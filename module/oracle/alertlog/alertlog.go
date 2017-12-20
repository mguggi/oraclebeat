package alertlog

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/paths"
	"github.com/elastic/beats/metricbeat/mb"
	"github.com/mguggi/oraclebeat/module/oracle"
	"github.com/peterbourgon/diskv"
	"github.com/pkg/errors"

	_ "gopkg.in/goracle.v2"
)

const selector = "oracle/alertlog"

// init registers the MetricSet with the central registry.
// The New method will be called after the setup of the module and before starting to fetch data
func init() {
	if err := mb.Registry.AddMetricSet("oracle", "alertlog", New, oracle.ParseDSN); err != nil {
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
	kvs     *diskv.Diskv
}

// New create a new instance of the MetricSet
// Part of new is also setting up the configuration by processing additional
// configuration entries if needed.
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {

	config := struct{}{}

	logp.Info("EXPERIMENTAL: The oracle alertlog metricset is experimental")

	if err := base.Module().UnpackConfig(&config); err != nil {
		return nil, err
	}

	// Configure the driver to return large floats as strings
	logp.Debug(selector, "Connect to host: %s", base.HostData().URI)
	db, err := sql.Open("goracle", base.HostData().URI)
	if err != nil {
		return nil, errors.Wrapf(err, "%s open connection failed", selector)
	}

	stmt, err := db.Prepare("select /* metricset: alertlog */ * from v$diag_alert_ext where originating_timestamp between :lastexec and :currexec")

	if err != nil {
		return nil, errors.Wrapf(err, "%s failed prepare statement", selector)
	}

	// Determine base path for persitent storage default is "path.data"
	basePath := paths.Paths.Data
	if len(basePath) == 0 {
		basePath = os.TempDir()
	}
	logp.Debug(selector, "The base path of persistent storage is '%s'", basePath)

	// Initialize state store, rooted at "path.data", with 8KB cache.
	kvs := diskv.New(diskv.Options{
		BasePath:     basePath,
		CacheSizeMax: 1024 * 8,
	})

	return &MetricSet{
		BaseMetricSet: base,
		kvs:           kvs,
		db:            db,
		stmt:          stmt,
	}, nil
}

// Fetch methods implements the data gathering and data conversion to the right format
// It returns the event which is then forward to the output. In case of an error, a
// descriptive error must be returned.
func (m *MetricSet) Fetch() ([]common.MapStr, error) {
	// Unique key for the key-value store
	key := oracle.GetHash(m.Host()) + ".oracle.alertlog.lastexec"

	// Make connection timeout aware
	timeout := m.BaseMetricSet.Module().Config().Timeout
	logp.Debug(selector, "Set timout to: %s", timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Set the current execution date
	currexec := time.Now()
	logp.Debug(selector, "Current execution date is [%s]", currexec)
	// Determine last execution date
	var lastexec time.Time
	if value, err := m.kvs.Read(key); err != nil {
		if os.IsNotExist(err) { // ignore because it is a normal behaviour on the first startup
			lastexec = time.Now().AddDate(-1, 0, 0) // default value on first startup
		} else {
			return nil, errors.Wrapf(err, "%s failed to load the last execution date", selector)
		}
	} else if len(value) > 0 {
		lastexec, err = time.Parse(time.RFC3339, string(value))
		if err != nil {
			return nil, errors.Wrapf(err, "%s failed to parse the last execution date", selector)
		}
	}
	logp.Debug(selector, "Last execution date is [%s]", lastexec)

	// Query the rows
	rows, err := m.stmt.QueryContext(ctx, lastexec, currexec)
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

	// Update last exectuion date
	if err := m.kvs.Write(key, []byte(currexec.Format(time.RFC3339Nano))); err != nil {
		return nil, errors.Wrapf(err, "%s failed to store the last execution date", selector)
	}

	return events, nil
}

// Closer is an optional interface that a MetricSet can implement in order to
// cleanup any resources it has open at shutdown.
func (m *MetricSet) Close() error {
	logp.Info("Cleanup the resources of %s", selector)

	var err error
	// close prepared statement
	if err = m.stmt.Close(); err != nil {
		errors.Wrapf(err, "%s failed close prepared statement", selector)
	}
	// close db connection
	if err = m.db.Close(); err != nil {
		errors.Wrapf(err, "%s failed close connection", selector)
	}

	return err
}
