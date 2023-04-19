package cloudwatch

import (
	"fmt"
	"grafana_to_datadog/grafana"
	"strings"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/gobeam/stringy"
)

type cloudwatchQuery struct {
	*grafana.Target
}

func (q *cloudwatchQuery) id() string {
	id := q.RefID
	if q.ID != "" {
		id = q.ID
	}
	return id
}

func (q *cloudwatchQuery) build() (string, error) {
	if q.QueryMode == "Logs" {
		return "", fmt.Errorf("unsupported query mode %s", q.QueryMode)
	}

	namespace := strings.Join(strings.Split(strings.ToLower(q.Namespace), "/"), ".")

	metricName := stringy.New(q.MetricName)
	metricName.SnakeCase().ToLower()
	metric := strings.Join([]string{namespace, metricName.SnakeCase().ToLower()}, ".")

	variables := []string{}
	for _, v := range q.Dimensions {
		if v != "*" {
			variables = append(variables, fmt.Sprint(v))
		}
	}

	from := "*"
	if len(variables) > 0 {
		from = strings.Join(variables, ",")
	}

	stats := "Average"
	if q.Statistic != "" {
		stats = q.Statistic
	} else if len(q.Statistics) > 0 {
		stats = q.Statistics[0]
	}

	return fmt.Sprintf("%s:%s{%s}.as_count()", statisticMap[stats], metric, from), nil
}

func (q *cloudwatchQuery) formula() (*datadogV1.WidgetFormula, error) {
	if q.Expression == "" {
		return datadogV1.NewWidgetFormula(q.id()), nil
	}
	formula := cloudwatchFormula{expression: q.Expression}
	return formula.build()
}
