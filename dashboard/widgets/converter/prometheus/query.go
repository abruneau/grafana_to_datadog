package prometheus

import (
	"fmt"
	"grafana_to_datadog/dashboard/widgets/shared"
	"grafana_to_datadog/dd"
	"strings"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

var aggregationMap = map[parser.ItemType]datadogV1.FormulaAndFunctionMetricAggregation{
	parser.AVG:   "avg",
	parser.SUM:   "sum",
	parser.MAX:   "max",
	parser.MIN:   "min",
	parser.COUNT: "sum",
	parser.SUB:   "-",
	parser.ADD:   "+",
	parser.MUL:   "*",
	parser.DIV:   "/",
}

type Query struct {
	*Target
	groupBy bool
	parsedExpr
}

type parsedExpr struct {
	agg     parser.ItemType
	groups  []string
	metric  string
	filters []*labels.Matcher
	err     error
}

func NewQuery(target map[string]interface{}, groupBy bool) (shared.Request, error) {
	var (
		t *Target
		q Query
	)
	t = shared.NewTarget[Target](target)
	q = Query{
		Target:  t,
		groupBy: groupBy,
	}
	return q.parseExpr()
}

func (q *Query) Build() (string, error) {
	var err error
	query := dd.Query{}

	if q.err != nil {
		return "", q.err
	}

	if q.metric == "" {
		return "", fmt.Errorf("no metric found query=%s", q.Expr)
	}
	query.Metric = q.metric

	query.Aggregator, err = q.Aggregator()
	if err != nil {
		return "", err
	}

	query.GroupBys = q.groups

	query.Function = q.function()

	query.Filters, err = q.filter()
	if err != nil {
		return "", err
	}

	return query.Build()
}

func (q *Query) Aggregator() (datadogV1.FormulaAndFunctionMetricAggregation, error) {
	var defaultValue datadogV1.FormulaAndFunctionMetricAggregation = "avg"

	if q.err != nil {
		return "", q.err
	}
	if q.agg == 0 {
		return defaultValue, nil
	}
	agg, ok := aggregationMap[q.agg]
	if !ok {
		return "", shared.AggregationTypeError(q.agg.String(), q.Expr)
	}
	return agg, nil
}

func (q *Query) function() dd.FormulaAndFunctionMetricFunction {
	return "as_count()"
}

func cleanupFilterValues(value string) []string {
	value = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(value, "..", "*"), ".*", "*"), "//", "/")
	values := strings.Split(value, "|")

	for i, v := range values {
		if strings.HasPrefix(v, "$") {
			values[i] = fmt.Sprintf("%s.value", v)
		}
	}
	return values
}

func filter(labelsMachers []*labels.Matcher) ([]string, error) {
	filters := []string{}

	for _, f := range labelsMachers {
		if f.Name == "__name__" {
			continue
		}

		values := cleanupFilterValues(f.Value)

		if len(values) > 1 {
			switch f.Type {
			case labels.MatchEqual:
				filters = append(filters, fmt.Sprintf("%s IN (%s)", f.Name, strings.Join(values, ", ")))
			case labels.MatchNotEqual:
				filters = append(filters, fmt.Sprintf("%s NOT IN (%s)", f.Name, strings.Join(values, ", ")))
			case labels.MatchRegexp, labels.MatchNotRegexp:
				return filters, fmt.Errorf("regex not supported with syntax operators \"IN\" and \"NOT IN\"")
			}
		} else {
			switch f.Type {
			case labels.MatchEqual, labels.MatchRegexp:
				filters = append(filters, fmt.Sprintf("%s:%s", f.Name, values[0]))
			case labels.MatchNotEqual, labels.MatchNotRegexp:
				filters = append(filters, fmt.Sprintf("!%s:%s", f.Name, values[0]))
			}
		}

	}

	return filters, nil
}

func (q *Query) filter() ([]string, error) {
	filters := []string{}

	for _, f := range q.filters {
		if f.Name == "__name__" {
			continue
		}

		values := cleanupFilterValues(f.Value)

		if len(values) > 1 {
			switch f.Type {
			case labels.MatchEqual:
				filters = append(filters, fmt.Sprintf("%s IN (%s)", f.Name, strings.Join(values, ", ")))
			case labels.MatchNotEqual:
				filters = append(filters, fmt.Sprintf("%s NOT IN (%s)", f.Name, strings.Join(values, ", ")))
			case labels.MatchRegexp, labels.MatchNotRegexp:
				return filters, fmt.Errorf("regex not supported with syntax operators \"IN\" and \"NOT IN\" query=%s", q.Expr)
			}
		} else {
			switch f.Type {
			case labels.MatchEqual, labels.MatchRegexp:
				filters = append(filters, fmt.Sprintf("%s:%s", f.Name, values[0]))
			case labels.MatchNotEqual, labels.MatchNotRegexp:
				filters = append(filters, fmt.Sprintf("!%s:%s", f.Name, values[0]))
			}
		}

	}

	return filters, nil
}
