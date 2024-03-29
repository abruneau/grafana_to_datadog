package prometheus

import (
	"fmt"
	"grafana_to_datadog/dashboard/widgets/shared"
	"grafana_to_datadog/dd"
	"log"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

type Structure struct {
	Function dd.FormulaAndFunctionMetricTransformation
	Args     []Structure
	Number   string
	Groups   []string
	Metric   string
	Filters  []*labels.Matcher
	Agg      datadogV1.FormulaAndFunctionMetricAggregation
	Parsed   string
}

var transfromationMap = map[string]dd.FormulaAndFunctionMetricTransformation{
	"abs":                "abs",
	"clamp_min":          "clamp_min",
	"clamp_max":          "clamp_max",
	"deriv":              "derivative",
	"log2":               "log2",
	"log10":              "log10",
	"delta":              "dt",
	"rate":               "per_second",
	"irate":              "per_second",
	"histogram_quantile": "histogram_quantile",
}

func extractAggregateFunction(expr parser.Expr) (dd.FormulaAndFunctionMetricFunction, parser.Expr, error) {

	f, ok := expr.(*parser.Call)
	if ok {
		if f.Func.Name == "rate" || f.Func.Name == "irate" {
			if len(f.Args) != 1 {
				return "", nil, fmt.Errorf("got more than 1 argument for rate function %s", f.String())
			}
			matrix, ok := f.Args[0].(*parser.MatrixSelector)
			if ok {
				return dd.FORMULAANDFUNCTIONMETRICFUNCTION_RATE, matrix.VectorSelector, nil
			}
		} else {
			return "", nil, fmt.Errorf("invalid function in aggregation %s %s", f.Func.Name, f.String())
		}
	}
	return dd.FORMULAANDFUNCTIONMETRICFUNCTION_COUNT, expr, nil
}

func parseMetric(expr parser.Expr) (name string, filters []string, err error) {
	vec, ok := expr.(*parser.VectorSelector)
	if ok {
		filters, err = filter(vec.LabelMatchers)
		return vec.Name, filters, err
	}
	return name, filters, fmt.Errorf("invalid expr type %s", expr.Type())
}

func parseAggregateExpr(expr parser.AggregateExpr, quantil float64) (agg datadogV1.FormulaAndFunctionMetricAggregation, query string, err error) {
	q := dd.Query{}
	var ok bool
	q.Aggregator, ok = aggregationMap[expr.Op]
	if !ok {
		return "", "", shared.AggregationTypeError(expr.Op.String(), expr.String())
	}

	if expr.Op == parser.QUANTILE {
		q.Aggregation = fmt.Sprintf("p%v", quantil*100)
	}

	q.GroupBys = expr.Grouping

	var metricExpr parser.Expr
	q.Function, metricExpr, err = extractAggregateFunction(expr.Expr)
	if err != nil {
		return "", "", err
	}
	q.Metric, q.Filters, err = parseMetric(metricExpr)
	if err != nil {
		return "", "", err
	}

	query, err = q.Build()
	return q.Aggregator, query, err
}

func (q *Query) parseVectorExpr(expr parser.VectorSelector) (s Structure, err error) {
	q.metric = expr.Name
	q.filters = expr.LabelMatchers
	s.Metric = expr.Name
	s.Filters = expr.LabelMatchers

	s.Parsed, err = q.Build()
	if err != nil {
		return s, err
	}
	s.Agg, err = q.Aggregator()
	if err != nil {
		return s, err
	}

	offset := expr.Offset + expr.OriginalOffset

	if offset > 0 {
		if offset <= 60*time.Minute {
			s.Function = dd.FORMULAANDFUNCTIONMETRICTRANSFORMATION_HOUR_BEFORE
		} else if offset <= 24*time.Hour {
			s.Function = dd.FORMULAANDFUNCTIONMETRICTRANSFORMATION_DAY_BEFORE
		} else if offset <= 7*24*time.Hour {
			s.Function = dd.FORMULAANDFUNCTIONMETRICTRANSFORMATION_WEEK_BEFORE
		} else {
			s.Function = dd.FORMULAANDFUNCTIONMETRICTRANSFORMATION_MONTH_BEFORE
		}
	}

	return
}

func handleQuantil(f *parser.Call) (s Structure, err error) {
	quantil := f.Args[0].(*parser.NumberLiteral).Val
	agg := f.Args[1].(*parser.AggregateExpr)
	agg.Op = parser.QUANTILE
	s.Agg, s.Parsed, err = parseAggregateExpr(*agg, quantil)
	return
}

func (q *Query) parseExprTypes(expr parser.Expr) (s Structure, err error) {
	num, ok := expr.(*parser.NumberLiteral)
	if ok {
		s.Number = num.String()
		return s, nil
	}

	agg, ok := expr.(*parser.AggregateExpr)
	if ok {
		s.Agg, s.Parsed, err = parseAggregateExpr(*agg, 0)
		if err != nil {
			return s, err
		}

		return s, nil
	}

	f, ok := expr.(*parser.Call)
	if ok {
		s.Function, ok = transfromationMap[f.Func.Name]
		if !ok {
			return s, fmt.Errorf("unsupported transformation function %s", f.Func.Name)
		}

		if s.Function == "histogram_quantile" {
			return handleQuantil(f)
		} else {

			for _, arg := range f.Args {
				if _, ok = arg.(*parser.NumberLiteral); ok {
					s.Args = append(s.Args, Structure{Number: arg.String()})
				} else {
					parsed, err := q.parseExprTypes(arg)
					if err != nil {
						return s, err
					}
					s.Args = append(s.Args, parsed)
					q.parseExprTypes(arg)
				}
			}
		}
		return s, nil
	}

	matrix, ok := expr.(*parser.MatrixSelector)
	if ok {
		parsed, err := q.parseExprTypes(matrix.VectorSelector)
		if err != nil {
			return s, err
		}
		s.Args = append(s.Args, parsed)
		return s, nil
	}

	vec, ok := expr.(*parser.VectorSelector)
	if ok {
		return q.parseVectorExpr(*vec)
	}

	bin, ok := expr.(*parser.BinaryExpr)
	if ok {
		lhs, err := q.parseExprTypes(bin.LHS)
		if err != nil {
			return s, err
		}
		rhs, err := q.parseExprTypes(bin.RHS)
		if err != nil {
			return s, err
		}

		s.Args = append(s.Args, lhs)

		s.Args = append(s.Args, rhs)
		s.Agg, ok = aggregationMap[bin.Op]
		if !ok {
			return s, shared.AggregationTypeError(bin.Op.String(), expr.String())
		}
	}
	return s, nil
}

func (q *Query) parseExpr() (r shared.Request, err error) {
	var expr parser.Expr

	q.groups = []string{}

	expr, err = parser.ParseExpr(q.Expr)
	if err != nil {
		fmt.Println(err)
		return r, fmt.Errorf("query parsing error: %s %v", q.Expr, err)
	}

	if expr.Type() != parser.ValueTypeVector {
		log.Fatalf(fmt.Errorf("expression type %s note supported", expr.Type()).Error())
		return r, fmt.Errorf("expression type %s note supported", expr.Type())
	}

	s, err := q.parseExprTypes(expr)
	if err != nil {
		return r, err
	}
	f, query, _ := s.transvers(q.RefID, 0)
	r.Formulas = append(r.Formulas, f)
	r.Queries = append(r.Queries, query...)
	return r, nil
}
