package cloudwatch

import (
	"testing"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/stretchr/testify/assert"
)

func TestMetric(t *testing.T) {
	var tests = []struct {
		input struct {
			namespace string
			queryName string
		}
		expected string
	}{
		{
			input: struct {
				namespace string
				queryName string
			}{"AWS/ApiGateway", "Count"},
			expected: "aws.apigateway.count",
		},
	}
	for _, test := range tests {
		testTarget := Query{&Target{}}
		testTarget.Namespace = test.input.namespace
		testTarget.MetricName = test.input.queryName
		m := testTarget.metric()
		assert.Equal(t, test.expected, m)
	}
}

func TestAggregator(t *testing.T) {
	for key, value := range statisticMap {
		testTarget := Query{&Target{}}
		testTarget.Statistic = key
		agg, err := testTarget.aggregator()
		assert.Equal(t, value, agg)
		assert.Nil(t, err)
	}
	for key, value := range statisticMap {
		testTarget := Query{&Target{}}
		testTarget.Statistics = []string{key}
		agg, err := testTarget.aggregator()
		assert.Equal(t, value, agg)
		assert.Nil(t, err)
	}

	testTarget := Query{&Target{}}
	agg, err := testTarget.aggregator()
	assert.Equal(t, datadogV1.FormulaAndFunctionMetricAggregation("avg"), agg)
	assert.Nil(t, err)
}
