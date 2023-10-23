package azure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetric(t *testing.T) {
	var tests = []struct {
		input struct {
			MetricName      string
			MetricNamespace string
		}
		expected string
	}{
		{
			input: struct {
				MetricName      string
				MetricNamespace string
			}{"Messages", "Microsoft.ServiceBus/namespaces"},
			expected: "azure.servicebus_namespaces.messages",
		},
	}
	for _, test := range tests {
		testTarget := Query{&Target{}, false}
		testTarget.AzureMonitor.MetricName = test.input.MetricName
		testTarget.AzureMonitor.MetricNamespace = test.input.MetricNamespace
		m, _ := testTarget.metric()
		assert.Equal(t, test.expected, m)
	}
}

func TestGroups(t *testing.T) {
	var tests = []struct {
		input []struct {
			Text  string `json:"text"`
			Value string `json:"value"`
		}
		expected []string
	}{
		{input: []struct {
			Text  string `json:"text"`
			Value string `json:"value"`
		}{
			{
				Text:  "None",
				Value: "None",
			},
			{
				Text:  "EntityName",
				Value: "EntityName",
			},
		},
			expected: []string{"entityname"},
		},
	}

	for _, test := range tests {
		testTarget := Query{&Target{}, true}
		testTarget.AzureMonitor.Dimensions = test.input
		m := testTarget.groups()
		assert.Equal(t, test.expected, m)
	}
}

func TestFilter(t *testing.T) {
	var tests = []struct {
		input []struct {
			Dimension string   `json:"dimension"`
			Operator  string   `json:"operator"`
			Filters   []string `json:"filters"`
		}
		expected []string
	}{
		{
			input: []struct {
				Dimension string   "json:\"dimension\""
				Operator  string   "json:\"operator\""
				Filters   []string "json:\"filters\""
			}{
				{
					Dimension: "EntityName",
					Operator:  "eq",
					Filters:   []string{"first"},
				},
			},
			expected: []string{"entityname:first"},
		},
		{
			input: []struct {
				Dimension string   "json:\"dimension\""
				Operator  string   "json:\"operator\""
				Filters   []string "json:\"filters\""
			}{
				{
					Dimension: "EntityName",
					Operator:  "ne",
					Filters:   []string{"first"},
				},
			},
			expected: []string{"!entityname:first"},
		},
		{
			input: []struct {
				Dimension string   "json:\"dimension\""
				Operator  string   "json:\"operator\""
				Filters   []string "json:\"filters\""
			}{
				{
					Dimension: "EntityName",
					Operator:  "sw",
					Filters:   []string{"first"},
				},
			},
			expected: []string{"entityname:first*"},
		},
		{
			input: []struct {
				Dimension string   "json:\"dimension\""
				Operator  string   "json:\"operator\""
				Filters   []string "json:\"filters\""
			}{
				{
					Dimension: "EntityName",
					Operator:  "eq",
					Filters:   []string{"first", "second"},
				},
			},
			expected: []string{"entityname IN (first, second)"},
		},
		{
			input: []struct {
				Dimension string   "json:\"dimension\""
				Operator  string   "json:\"operator\""
				Filters   []string "json:\"filters\""
			}{
				{
					Dimension: "EntityName",
					Operator:  "ne",
					Filters:   []string{"first", "second"},
				},
			},
			expected: []string{"entityname NOT IN (first, second)"},
		},
		{
			input: []struct {
				Dimension string   "json:\"dimension\""
				Operator  string   "json:\"operator\""
				Filters   []string "json:\"filters\""
			}{
				{
					Dimension: "EntityName",
					Operator:  "sw",
					Filters:   []string{"first", "second"},
				},
			},
			expected: []string{"entityname IN (first*, second*)"},
		},
	}

	for _, test := range tests {
		testTarget := Query{&Target{}, false}
		testTarget.AzureMonitor.DimensionFilters = test.input
		f, _ := testTarget.filter()
		assert.Equal(t, test.expected, f)
	}
}
