package cloudwatch

import (
	"fmt"
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
}

func (q *Query) id() string {
	id := q.RefID
	if q.ID != "" {
		id = q.ID
	}
	return id
}

func (q *Query) metric() string {
	namespace := strings.Join(strings.Split(strings.ToLower(q.Namespace), "/"), ".")

	metricName := stringy.New(q.MetricName)
	metricName.SnakeCase().ToLower()
	return strings.Join([]string{namespace, metricName.SnakeCase().ToLower()}, ".")
}

func (q *Query) filter() []string {
	variables := []string{}
	for _, v := range q.Dimensions {
		if v != "*" {
			variables = append(variables, fmt.Sprint(v))
		}
	}
	return variables
}

// TODO: to be implemented
func (q *Query) groups() []string {
	return []string{}
}

func (q *Query) aggregator() (datadogV1.FormulaAndFunctionMetricAggregation, error) {
	stats := "Average"
	if q.Statistic != "" {
		stats = q.Statistic
	} else if len(q.Statistics) > 0 {
		stats = q.Statistics[0]
	}
	agg, ok := statisticMap[stats]

	if !ok {
		return "", fmt.Errorf("alignement type %s not supported", q.MetricQuery.PerSeriesAligner)
	}
	return agg, nil
}

// TODO: It should not be static
func (q *Query) function() dd.FormulaAndFunctionMetricFunction {
	return "as_count()"
}

func (q *Query) build() (string, error) {
	var err error
	query := dd.Query{}
	if q.QueryMode == "Logs" {
		return "", fmt.Errorf("unsupported query mode %s", q.QueryMode)
	}

	query.Metric = q.metric()

	query.Filters = q.filter()

	query.GroupBys = q.groups()

	query.Aggregator, err = q.aggregator()
	if err != nil {
		return "", err
	}

	query.Function = q.function()

	return query.Build()
}

func (q *Query) formula() (*datadogV1.WidgetFormula, error) {
	if q.Expression == "" {
		return datadogV1.NewWidgetFormula(q.id()), nil
	}
	formula := cloudwatchFormula{expression: q.Expression}
	return formula.build()
}
