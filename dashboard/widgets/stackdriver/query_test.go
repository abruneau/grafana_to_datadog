package stackdriver

import (
	"grafana_to_datadog/dd"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetric(t *testing.T) {

	var tests = []struct {
		input    string
		expected string
	}{
		{
			"compute.googleapis.com/instance/uptime", "gcp.gce.instance.uptime",
		},
		{
			"compute.googleapis.com/instance/disk/read_ops_count", "gcp.gce.instance.disk.read_ops_count",
		},
		{
			"compute.googleapis.com/instance/disk/write_ops_count", "gcp.gce.instance.disk.write_ops_count",
		},
	}

	for _, test := range tests {
		testTarget := Query{&Target{}}
		testTarget.MetricQuery.MetricType = test.input
		m, _ := testTarget.metric()
		assert.Equal(t, test.expected, m)
	}

	testTarget := Query{&Target{}}
	testTarget.MetricQuery.Filters = []string{"resource.type",
		"=",
		"gce_instance",
		"AND",
		"metric.type",
		"=",
		"compute.googleapis.com/instance/cpu/utilization",
	}
	m, _ := testTarget.metric()
	assert.Equal(t, "gcp.gce.instance.cpu.utilization", m)

	testTarget = Query{&Target{}}
	testTarget.MetricQuery.Filters = []string{"metric.label.instance_name",
		"=",
		"gke-agent-core-benchmark-default-pool-4b1e15cf-bue8",
		"AND",
		"metric.type",
		"=",
		"compute.googleapis.com/guest/cpu/load_15m",
	}
	m, _ = testTarget.metric()
	assert.Equal(t, "gcp.gce.guest.cpu.load_15m", m)

}

func TestGroups(t *testing.T) {
	var tests = []struct {
		input    string
		expected []string
	}{
		{
			"{{metric.label.instance_name}}", []string{"instance_name"},
		},
		{
			"", []string{},
		},
	}

	for _, test := range tests {
		testTarget := Query{&Target{}}
		testTarget.MetricQuery.AliasBy = test.input
		m, _ := testTarget.groups()
		assert.Equal(t, test.expected, m)
	}
}

func TestAggregator(t *testing.T) {
	for key, value := range alignmentType {
		testTarget := Query{&Target{}}
		testTarget.MetricQuery.PerSeriesAligner = key
		agg, err := testTarget.aggregator()
		assert.Equal(t, value, agg)
		assert.Nil(t, err)
	}
}

func TestFilter(t *testing.T) {
	var tests = []struct {
		input    []string
		expected []string
	}{
		{
			[]string{"metric.label.instance_name", "=", "gke-agent-core-benchmark-default-pool-4b1e15cf-gqiy", "AND", "metric.type", "=", "compute.googleapis.com/guest/cpu/load_15m"},
			[]string{"instance_name:gke-agent-core-benchmark-default-pool-4b1e15cf-gqiy"},
		},
	}

	for _, test := range tests {
		testTarget := Query{&Target{}}
		testTarget.TimeSeriesList.Filters = test.input
		testTarget.QueryType = "timeSeriesList"
		m, _ := testTarget.filter()
		assert.Equal(t, test.expected, m)
	}
}

func TestFunction(t *testing.T) {
	testTarget := Query{&Target{}}
	testTarget.MetricQuery.PerSeriesAligner = "ALIGN_RATE"
	f := testTarget.function()
	assert.Equal(t, dd.FORMULAANDFUNCTIONMETRICFUNCTION_RATE, f)
	testTarget = Query{&Target{}}
	testTarget.MetricQuery.PerSeriesAligner = "ALIGN_DELTA"
	f = testTarget.function()
	assert.Equal(t, dd.FORMULAANDFUNCTIONMETRICFUNCTION_COUNT, f)
	testTarget = Query{&Target{}}
	testTarget.MetricQuery.PerSeriesAligner = ""
	testTarget.MetricQuery.MetricKind = "DELTA"
	f = testTarget.function()
	assert.Equal(t, dd.FORMULAANDFUNCTIONMETRICFUNCTION_COUNT, f)
	testTarget = Query{&Target{}}
	testTarget.MetricQuery.PerSeriesAligner = ""
	testTarget.MetricQuery.MetricKind = ""
	f = testTarget.function()
	assert.Equal(t, dd.FormulaAndFunctionMetricFunction(""), f)
}
