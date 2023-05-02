package stackdriver

import (
	"fmt"
	"grafana_to_datadog/dd"
	"strings"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"golang.org/x/exp/slices"
)

var (
	root = "gcp"
)

var alignmentType = map[string]datadogV1.FormulaAndFunctionMetricAggregation{
	"ALIGN_MEAN":  "avg",
	"ALIGN_MIN":   "min",
	"ALIGN_MAX":   "max",
	"ALIGN_SUM":   "sum",
	"ALIGN_RATE":  "avg",
	"ALIGN_DELTA": "avg",
}

type Query struct {
	*Target
}

func (q *Query) id() string {
	id := q.RefID
	// if q.ID != "" {
	// 	id = q.ID
	// }
	return id
}

func resourceMap(res string) string {
	switch res {
	case "compute.googleapis.com":
		return "gce"
	default:
		return strings.Split(res, ".")[0]
	}
}

func (q *Query) metric() (string, error) {

	m := q.MetricQuery.MetricType
	if m == "" {
		idx := slices.Index(q.MetricQuery.Filters, "metric.type")
		if idx != -1 && len(q.MetricQuery.Filters) > idx+2 {
			m = q.MetricQuery.Filters[idx+2]
		} else {
			return "", fmt.Errorf("no metric found")
		}
	}

	parts := strings.Split(m, "/")
	resourceType := resourceMap(parts[0])
	metric := strings.Join(parts[1:], ".")

	return strings.Join([]string{root, resourceType, metric}, "."), nil
}

func (q *Query) filter() ([]string, error) {
	filters := []string{}
	if q.MetricQuery.ProjectName != "" {
		filters = append(filters, q.MetricQuery.ProjectName)
	}
	return filters, nil
}

func (q *Query) groups() ([]string, error) {
	groupBys := []string{}

	if q.MetricQuery.AliasBy != "" {
		if !(strings.HasPrefix(q.MetricQuery.AliasBy, "{{") && strings.HasSuffix(q.MetricQuery.AliasBy, "}}")) {
			return groupBys, fmt.Errorf("unsuported alias %s", q.MetricQuery.AliasBy)
		}
		alias := q.MetricQuery.AliasBy[2 : len(q.MetricQuery.AliasBy)-2]
		splitAlias := strings.Split(alias, ".")[2:]
		groupBys = append(groupBys, strings.Join(splitAlias, "."))
	}

	for _, g := range q.MetricQuery.GroupBys {
		groupBys = append(groupBys, strings.Join(strings.Split(g, ".")[2:], "."))
	}

	return groupBys, nil
}

func (q *Query) aggregator() (datadogV1.FormulaAndFunctionMetricAggregation, error) {
	agg, ok := alignmentType[q.MetricQuery.PerSeriesAligner]
	if !ok {
		return "", fmt.Errorf("alignement type %s not supported", q.MetricQuery.PerSeriesAligner)
	}
	return agg, nil
}

func (q *Query) function() dd.FormulaAndFunctionMetricFunction {
	if q.MetricQuery.PerSeriesAligner == "ALIGN_RATE" {
		return "as_rate()"
	}

	if q.MetricQuery.PerSeriesAligner == "ALIGN_DELTA" {
		return "as_count()"
	}

	if q.MetricQuery.MetricKind == "DELTA" {
		return "as_count()"
	}

	return ""
}

func (q *Query) build() (string, error) {
	var err error
	query := dd.Query{}

	if q.QueryType != "metrics" {
		return "", fmt.Errorf("unsupported query mode %s", q.QueryType)
	}

	query.Metric, err = q.metric()
	if err != nil {
		return "", err
	}

	query.Aggregator, err = q.aggregator()
	if err != nil {
		return "", err
	}

	query.GroupBys, err = q.groups()
	if err != nil {
		return "", err
	}

	query.Function = q.function()

	query.Filters, err = q.filter()
	if err != nil {
		return "", err
	}

	return query.Build()
}

func (q *Query) formula() string {
	if q.Hide {
		return ""
	}
	return q.id()
}
