package dd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	input    Query
	expected string
	err      error
}{
	{
		input:    Query{Metric: "system.cpu.user", Aggregator: "avg"},
		expected: "avg:system.cpu.user{*}",
		err:      nil,
	},
	{
		input:    Query{Metric: "system.cpu.user", Aggregator: "avg", Filters: []string{"host:a", "user:a"}},
		expected: "avg:system.cpu.user{host:a,user:a}",
		err:      nil,
	},
	{
		input:    Query{Metric: "system.cpu.user", Aggregator: "avg", Filters: []string{"host:a", "user:a"}, GroupBys: []string{"host", "user"}},
		expected: "avg:system.cpu.user{host:a,user:a} by {host,user}",
		err:      nil,
	},
	{
		input:    Query{Metric: "system.cpu.user", Aggregator: "sum"},
		expected: "sum:system.cpu.user{*}",
		err:      nil,
	},
	{
		input:    Query{Metric: "system.cpu.user", Aggregator: "test"},
		expected: "",
		err:      fmt.Errorf("unknown agrregator test"),
	},
	{
		input:    Query{Metric: "system.cpu.user", Aggregator: "sum", Function: "as_count()"},
		expected: "sum:system.cpu.user{*}.as_count()",
		err:      nil,
	},
	{
		input:    Query{Metric: "system.cpu.user", Aggregator: "sum", Function: "as_test()"},
		expected: "",
		err:      fmt.Errorf("unknown function as_test()"),
	},
}

func TestQueryBuild(t *testing.T) {
	for _, test := range tests {
		q, err := test.input.Build()
		assert.Equal(t, test.expected, q)
		assert.Equal(t, test.err, err)
	}
}
