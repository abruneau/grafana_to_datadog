package prometheus

import (
	"fmt"
	"grafana_to_datadog/dashboard/widgets/shared"
	"testing"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	input         string
	output        parsedExpr
	filters       []string
	filters_error error
}{
	{
		input: "sum(istio_build{component=\"pilot\"}) by (tag)",
		output: parsedExpr{
			agg:     parser.SUM,
			groups:  []string{"tag"},
			metric:  "istio_build",
			filters: []*labels.Matcher{labels.MustNewMatcher(labels.MatchEqual, "__name__", "istio_build"), labels.MustNewMatcher(labels.MatchEqual, "component", "pilot")},
		},
		filters: []string{"component:pilot"},
	},
	{
		input: "process_virtual_memory_bytes{app=\"istiod\"}",
		output: parsedExpr{
			groups:  []string{},
			metric:  "process_virtual_memory_bytes",
			filters: []*labels.Matcher{labels.MustNewMatcher(labels.MatchEqual, "__name__", "process_virtual_memory_bytes"), labels.MustNewMatcher(labels.MatchEqual, "app", "istiod")},
		},
		filters: []string{"app:istiod"},
	},
	{
		input: "container_memory_working_set_bytes{container=~\"istio-proxy\", pod=~\"istiod-.*|istio-pilot-.*\"}",
		output: parsedExpr{
			groups:  []string{},
			metric:  "container_memory_working_set_bytes",
			filters: []*labels.Matcher{labels.MustNewMatcher(labels.MatchEqual, "__name__", "container_memory_working_set_bytes"), labels.MustNewMatcher(labels.MatchRegexp, "container", "istio-proxy"), labels.MustNewMatcher(labels.MatchRegexp, "pod", "istiod-.*|istio-pilot-.*")},
		},
		filters:       []string{"container:istio-proxy", "pod IN (istiod-*, istio-pilot-*)"},
		filters_error: fmt.Errorf("regex not supported with syntax operators \"IN\" and \"NOT IN\" query=%s", "container_memory_working_set_bytes{container=~\"istio-proxy\", pod=~\"istiod-.*|istio-pilot-.*\"}"),
	},
	{
		input: "sum(irate(container_cpu_usage_seconds_total{container=\"discovery\", pod=~\"istiod-.*|istio-pilot-.*\"}[1m]))",
		output: parsedExpr{
			agg:       parser.SUM,
			groups:    []string{},
			functions: []metricFunction{{name: "irate"}},
			metric:    "container_cpu_usage_seconds_total",
			filters:   []*labels.Matcher{labels.MustNewMatcher(labels.MatchEqual, "__name__", "container_cpu_usage_seconds_total"), labels.MustNewMatcher(labels.MatchRegexp, "container", "discovery"), labels.MustNewMatcher(labels.MatchRegexp, "pod", "istiod-.*|istio-pilot-.*")},
		},
		filters:       []string{"container:discovery", "pod IN (istiod-*, istio-pilot-*)"},
		filters_error: fmt.Errorf("regex not supported with syntax operators \"IN\" and \"NOT IN\" query=%s", "sum(irate(container_cpu_usage_seconds_total{container=\"discovery\", pod=~\"istiod-.*|istio-pilot-.*\"}[1m]))"),
	},
}

func TestParseExpr(t *testing.T) {
	for _, test := range tests {
		testTarget := Query{&Target{}, false, parsedExpr{}}
		testTarget.Expr = test.input
		testTarget.err = testTarget.parseExpr()
		assert.NoError(t, testTarget.err)
		assert.Equal(t, test.output.agg, testTarget.parsedExpr.agg)
		assert.Equal(t, test.output.groups, testTarget.parsedExpr.groups)
		assert.Equal(t, test.output.functions, testTarget.parsedExpr.functions)
		assert.Equal(t, test.output.metric, testTarget.parsedExpr.metric)
		assert.ObjectsAreEqualValues(test.output.filters, testTarget.parsedExpr.filters)
	}
}

func TestFilter(t *testing.T) {
	for _, test := range tests {
		testTarget := Query{&Target{}, false, parsedExpr{}}
		testTarget.Expr = test.input
		testTarget.err = testTarget.parseExpr()
		assert.NoError(t, testTarget.err)
		filters, err := testTarget.filter()
		if test.filters_error != nil {
			assert.EqualError(t, err, test.filters_error.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.filters, filters)
		}

	}
}

func TestCleanupFilterValues(t *testing.T) {
	var examples = []struct {
		value  string
		values []string
	}{
		{
			value:  "istiod-.*|istio-pilot-.*",
			values: []string{"istiod-*", "istio-pilot-*"},
		},
		{
			value:  "$cluster",
			values: []string{"$cluster.value"},
		},
	}

	for _, e := range examples {
		assert.Equal(t, e.values, cleanupFilterValues(e.value))
	}
}

func TestAggregator(t *testing.T) {

	// Test all values
	for key, value := range agregationMap {
		testTarget := Query{&Target{}, false, parsedExpr{}}
		testTarget.agg = key
		agg, err := testTarget.Aggregator()
		assert.Equal(t, value, agg)
		assert.Nil(t, err)
	}

}

func TestAggregatorDefault(t *testing.T) {
	testTarget := Query{&Target{}, false, parsedExpr{}}
	agg, err := testTarget.Aggregator()
	assert.Equal(t, datadogV1.FormulaAndFunctionMetricAggregation("avg"), agg)
	assert.Nil(t, err)
}

func TestAggregatorError(t *testing.T) {
	testTarget := Query{&Target{}, false, parsedExpr{}}
	testTarget.agg = parser.BOTTOMK
	_, err := testTarget.Aggregator()
	assert.EqualError(t, err, shared.AggregationTypeError("bottomk").Error())
}
