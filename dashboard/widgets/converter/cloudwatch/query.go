package cloudwatch

import (
	"fmt"
	"grafana_to_datadog/dashboard/widgets/shared"
	"grafana_to_datadog/dd"
	"strings"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/gobeam/stringy"
)

var statisticMap = map[string]datadogV1.FormulaAndFunctionMetricAggregation{
	"Average": "avg",
	"Sum":     "sum",
	"Maximum": "max",
}

type Query struct {
	*Target
	groupBy bool
}

func NewQuery(target map[string]interface{}, groupBy bool) (shared.Request, error) {
	t := shared.NewTarget[Target](target)
	query := &Query{
		Target:  t,
		groupBy: groupBy,
	}
	return query.parse()
}

func (q *Query) parse() (r shared.Request, err error) {
	query, err := q.Build()
	if err != nil {
		return r, err
	}
	agg, err := q.Aggregator()
	if err != nil {
		return r, err
	}
	r.Queries = append(r.Queries, struct {
		Name        string
		Query       string
		Aggregation datadogV1.FormulaAndFunctionMetricAggregation
	}{q.Id(), query, agg})
	r.Formulas = append(r.Formulas, q.Formula())
	return r, nil
}

func (q *Query) Id() string {
	id := q.RefID
	if q.ID != "" {
		id = q.ID
	}
	return id
}

func (q *Query) metric() string {
	namespace := strings.Join(strings.Split(strings.ToLower(q.Namespace), "/"), ".")

	metricName := stringy.New(q.MetricName)
	return strings.Join([]string{namespace, metricName.SnakeCase().ToLower()}, ".")
}

func (q *Query) filter() []string {
	variables := []string{}
	for dim, v := range q.Dimensions {
		if v != "*" {
			dimName := stringy.New(dim)
			variables = append(variables, fmt.Sprintf("%s:%s", dimName.SnakeCase().ToLower(), v))
		}
	}
	if q.Region != "" {
		variables = append(variables, q.Region)
	}
	return variables
}

func (q *Query) groups() []string {
	variables := []string{}
	if !q.groupBy {
		return variables
	}
	for dim, v := range q.Dimensions {
		if v == "*" {
			dimName := stringy.New(dim)
			variables = append(variables, dimName.SnakeCase().ToLower())
		}
	}
	return variables
}

func (q *Query) Aggregator() (datadogV1.FormulaAndFunctionMetricAggregation, error) {
	stats := "Average"
	if q.Statistic != "" {
		stats = q.Statistic
	} else if len(q.Statistics) > 0 {
		stats = q.Statistics[0]
	}
	agg, ok := statisticMap[stats]

	if !ok {
		return "", shared.AggregationTypeError(q.MetricQuery.PerSeriesAligner)
	}
	return agg, nil
}

// TODO: It should not be static
func (q *Query) function() dd.FormulaAndFunctionMetricFunction {
	return "as_count()"
}

func (q *Query) Build() (string, error) {
	var err error
	query := dd.Query{}
	if q.QueryMode == "Logs" {
		return "", fmt.Errorf("unsupported query mode %s", q.QueryMode)
	}

	query.Metric = q.metric()

	query.Filters = q.filter()

	query.GroupBys = q.groups()

	query.Aggregator, err = q.Aggregator()
	if err != nil {
		return "", err
	}

	query.Function = q.function()

	return query.Build()
}

func (q *Query) Formula() string {
	if q.Type != "math" {
		return q.Id()
	}

	if q.Expression == "" {
		return ""
	}

	return strings.ReplaceAll(q.Expression, "$", "")
}
