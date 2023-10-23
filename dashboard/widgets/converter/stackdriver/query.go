package stackdriver

import (
	"fmt"
	"grafana_to_datadog/dashboard/widgets/shared"
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
	groupBy bool
}

func NewQuery(target map[string]interface{}, groupBy bool) shared.Query {
	t := shared.NewTarget[Target](target)
	query := &Query{
		Target:  t,
		groupBy: groupBy,
	}
	return query
}

func (q *Query) Id() string {
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
		var queryFilters []string
		if q.QueryType == "timeSeriesList" {
			queryFilters = q.TimeSeriesList.Filters
		} else {
			queryFilters = q.MetricQuery.Filters
		}

		idx := slices.Index(queryFilters, "metric.type")
		if idx != -1 && len(queryFilters) > idx+2 {
			m = queryFilters[idx+2]
		}

		if m == "" {
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
	var queryFilters []string
	if q.QueryType == "timeSeriesList" {
		queryFilters = q.TimeSeriesList.Filters
	} else {
		queryFilters = q.MetricQuery.Filters
	}

	for i, f := range queryFilters {
		if strings.HasPrefix(f, "metric.label") {
			label := strings.Split(f, ".")[2]
			value := queryFilters[i+2]
			filters = append(filters, fmt.Sprintf("%s:%s", label, value))
		}
	}

	return filters, nil
}

func (q *Query) groups() ([]string, error) {
	groupBys := []string{}
	if !q.groupBy {
		return groupBys, nil
	}

	if q.QueryType == "timeSeriesList" && q.TimeSeriesList.AliasBy != "" {
		if !(strings.HasPrefix(q.TimeSeriesList.AliasBy, "{{") && strings.HasSuffix(q.TimeSeriesList.AliasBy, "}}")) {
			return groupBys, fmt.Errorf("unsuported alias %s", q.TimeSeriesList.AliasBy)
		}
		alias := q.TimeSeriesList.AliasBy[2 : len(q.TimeSeriesList.AliasBy)-2]
		splitAlias := strings.Split(alias, ".")[2:]
		groupBys = append(groupBys, strings.Join(splitAlias, "."))
	} else if q.MetricQuery.AliasBy != "" {
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

	for _, g := range q.TimeSeriesList.GroupBys {
		groupBys = append(groupBys, strings.Join(strings.Split(g, ".")[2:], "."))
	}

	return groupBys, nil
}

func (q *Query) Aggregator() (datadogV1.FormulaAndFunctionMetricAggregation, error) {
	var aligner string
	if q.QueryType == "timeSeriesList" {
		aligner = q.TimeSeriesList.PerSeriesAligner
	} else {
		aligner = q.MetricQuery.PerSeriesAligner
	}
	agg, ok := alignmentType[aligner]
	if !ok {
		return "", shared.AggregationTypeError(aligner)
	}
	return agg, nil
}

func (q *Query) function() dd.FormulaAndFunctionMetricFunction {
	var aligner string
	var metricKind string
	if q.QueryType == "timeSeriesList" {
		aligner = q.TimeSeriesList.PerSeriesAligner
		metricKind = q.TimeSeriesList.MetricKind
	} else {
		aligner = q.MetricQuery.PerSeriesAligner
		metricKind = q.MetricQuery.MetricKind
	}
	if aligner == "ALIGN_RATE" {
		return "as_rate()"
	}

	if aligner == "ALIGN_DELTA" {
		return "as_count()"
	}

	if metricKind == "DELTA" {
		return "as_count()"
	}

	return ""
}

func (q *Query) Build() (string, error) {
	var err error
	query := dd.Query{}

	if !slices.Contains([]string{"metrics", "timeSeriesList"}, q.QueryType) {
		return "", fmt.Errorf("unsupported query mode %s", q.QueryType)
	}

	query.Metric, err = q.metric()
	if err != nil {
		return "", err
	}

	query.Aggregator, err = q.Aggregator()
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

func (q *Query) Formula() string {
	if q.Hide {
		return ""
	}
	return q.Id()
}
