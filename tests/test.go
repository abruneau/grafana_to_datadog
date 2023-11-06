package main

import (
	"fmt"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

func parseExp(expr parser.Expr) {
	fmt.Println("Parsing expr of type ", expr.Type())
	agg, ok := expr.(*parser.AggregateExpr)
	if ok {
		fmt.Println("Parsing AggregateExpr")
		fmt.Printf("Agg: %s, Groups: %v\n", agg.Op, agg.Grouping)
		parseExp(agg.Expr)
	}

	f, ok := expr.(*parser.Call)
	if ok {
		fmt.Println("Parsing Call")
		fmt.Printf("Function: %s\n", f.Func.Name)
		fmt.Println("Function args: ", len(f.Args))
		for _, arg := range f.Args {
			fmt.Println(arg)
			parseExp(arg)
		}
	}

	vec, ok := expr.(*parser.VectorSelector)
	if ok {
		fmt.Println("Parsing VectorSelector")
		fmt.Println("Metric: ", vec.Name)
		fmt.Println("Fileters: ", vec.LabelMatchers)
	}

	matrix, ok := expr.(*parser.MatrixSelector)
	if ok {
		parseExp(matrix.VectorSelector)
	}
}

type Query struct {
	Expr string
	parsedExpr
}

type parsedExpr struct {
	agg       parser.ItemType
	groups    []string
	metric    string
	filters   []*labels.Matcher
	functions []metricFunction
}

type metricFunction struct {
	name  string
	value []any
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
				mf.value = append(mf.value, arg.String())
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

func main() {

	query := "sum(pilot_xds_eds_reject{app=\"istiod\"}) or (absent(pilot_xds_eds_reject{app=\"istiod\"}) - 1)"

	q := Query{Expr: query}

	q.parseExpr()

	fmt.Println("metric: ", q.metric)
	fmt.Println("agg: ", q.agg)
	fmt.Println("groups: ", q.groups)
	fmt.Println("filters: ", q.filters)
	fmt.Println("functions: ", q.functions)

}
