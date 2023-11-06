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

var agregationMap = map[parser.ItemType]datadogV1.FormulaAndFunctionMetricAggregation{
	parser.AVG:   "avg",
	parser.SUM:   "sum",
	parser.MAX:   "max",
	parser.MIN:   "min",
	parser.COUNT: "sum",
}

// var transfromationMap = map[string]dd.FormulaAndFunctionMetricTransformation{
// 	"abs":       "abs",
// 	"clamp_min": "clamp_min",
// 	"clamp_max": "clamp_max",
// 	"deriv":     "derivative",
// 	"log2":      "log2",
// 	"log10":     "log10",
// 	"delta":     "dt",
// 	"rate":      "per_second",
// 	"irate":     "per_second",
// }

type Query struct {
	*Target
	groupBy bool
	parsedExpr
}

type parsedExpr struct {
	agg       parser.ItemType
	groups    []string
	metric    string
	filters   []*labels.Matcher
	functions []metricFunction
	err       error
}

type metricFunction struct {
	name   string
	values []string
}

func NewQuery(target map[string]interface{}, groupBy bool) shared.Query {
	t := shared.NewTarget[Target](target)
	query := &Query{
		Target:  t,
		groupBy: groupBy,
	}
	query.err = query.parseExpr()
	return query
}

func (q *Query) Id() string {
	id := q.RefID
	return id
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

	// query.Transformations, err = q.transfromations()
	// if err != nil {
	// 	return "", err
	// }

	return query.Build()
}

func (q *Query) parseExprTypes(expr parser.Expr) {
	agg, ok := expr.(*parser.AggregateExpr)
	if ok {
		q.agg = agg.Op
		if agg.Grouping != nil {
			q.groups = agg.Grouping
		}
		q.parseExprTypes(agg.Expr)
	}

	f, ok := expr.(*parser.Call)
	if ok {
		mf := metricFunction{name: f.Func.Name}

		for _, arg := range f.Args {
			if _, ok = arg.(*parser.NumberLiteral); ok {
				mf.values = append(mf.values, arg.String())
			} else {
				q.parseExprTypes(arg)
			}
		}
		q.functions = append(q.functions, mf)
	}

	matrix, ok := expr.(*parser.MatrixSelector)
	if ok {
		q.parseExprTypes(matrix.VectorSelector)
	}

	vec, ok := expr.(*parser.VectorSelector)
	if ok {
		q.metric = vec.Name
		q.filters = vec.LabelMatchers
	}
}

func (q *Query) parseExpr() error {
	var expr parser.Expr
	var err error

	q.groups = []string{}

	expr, err = parser.ParseExpr(q.Expr)
	if err != nil {
		return fmt.Errorf("query parsing error: %s %v", q.Expr, err)
	}

	if expr.Type() != parser.ValueTypeVector {
		return fmt.Errorf("expression type %s note supported", expr.Type())
	}

	q.parseExprTypes(expr)
	return nil
}

func (q *Query) Aggregator() (datadogV1.FormulaAndFunctionMetricAggregation, error) {
	var defaultValue datadogV1.FormulaAndFunctionMetricAggregation = "avg"

	if q.err != nil {
		return "", q.err
	}
	if q.agg == 0 {
		return defaultValue, nil
	}
	agg, ok := agregationMap[q.agg]
	if !ok {
		return "", shared.AggregationTypeError(q.agg.String(), q.Expr)
	}
	return agg, nil
}

func (q *Query) Formula() string {
	if q.Hide {
		return ""
	}
	return q.Id()
}

func (q *Query) function() dd.FormulaAndFunctionMetricFunction {
	return "as_count()"
}

func (q *Query) filter() ([]string, error) {
	filters := []string{}

	for _, f := range q.filters {
		if f.Name == "__name__" {
			continue
		}

		value := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(f.Value, "..", "*"), ".*", "*"), "//", "/")
		values := strings.Split(value, "|")

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
				filters = append(filters, fmt.Sprintf("%s:%s", f.Name, value))
			case labels.MatchNotEqual, labels.MatchNotRegexp:
				filters = append(filters, fmt.Sprintf("!%s:%s", f.Name, value))
			}
		}

	}

	return filters, nil
}

// func (q *Query) transfromations() ([]dd.MetricFunction, error) {
// 	var transfo []dd.MetricFunction
// 	for _, f := range q.functions {
// 		var ok bool
// 		mf := dd.MetricFunction{Values: f.values}
// 		mf.Name, ok = transfromationMap[f.name]
// 		if !ok {
// 			return transfo, shared.TransformationTypeError(f.name, q.Expr)
// 		}
// 		transfo = append(transfo, mf)
// 	}
// 	return transfo, nil
// }
