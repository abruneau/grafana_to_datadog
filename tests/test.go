package main

import (
	"fmt"
	"log"

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

func main() {

	query := "sum(irate(container_cpu_usage_seconds_total{container=\"discovery\", pod=~\"istiod-.*|istio-pilot-.*\"}[1m]))"

	var expr parser.Expr
	var err error

	expr, err = parser.ParseExpr(query)
	if err != nil {
		log.Fatalf("parse error: %s", err)
	}

	parseExp(expr)

}
