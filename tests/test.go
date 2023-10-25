package main

import (
	"fmt"
	"log"

	"github.com/VictoriaMetrics/metricsql"
	"github.com/prometheus/prometheus/promql/parser"
)

func metricsqlTest(query string) {
	expr, err := metricsql.Parse(query)
	if err != nil {
		log.Fatalf("parse error: %s", err)
	}
	fmt.Printf("parsed expr: %s\n", expr.AppendString(nil))

	ae := expr.(*metricsql.AggrFuncExpr)
	fmt.Printf("aggr func: name=%s, modifier=%s\n", ae.Name, ae.Modifier.AppendString(nil))

	me := ae.Args[0].(*metricsql.MetricExpr)
	for _, v := range me.LabelFilterss[0] {
		fmt.Printf("name: %s, value: %s\n", v.Label, v.Value)
	}
}

func main() {

	query := "sum(istio_build{component=\"pilot\"}) by (tag)"

	var expr parser.Expr
	var err error

	expr, err = parser.ParseExpr(query)
	if err != nil {
		log.Fatalf("parse error: %s", err)
	}

	fmt.Println(expr.Type())

	agg, ok := expr.(*parser.AggregateExpr)
	if ok {
		fmt.Printf("Agg: %s, Groups: %v\n", agg.Op, agg.Grouping)
		expr = agg.Expr
	}

	vec, ok := expr.(*parser.VectorSelector)
	if ok {
		fmt.Println("Metric: ", vec.Name)
		fmt.Println("Fileters: ", vec.LabelMatchers)
	}

}
