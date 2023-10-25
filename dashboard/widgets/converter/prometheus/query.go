package prometheus

import (
	"fmt"
	"grafana_to_datadog/dashboard/widgets/shared"
	"grafana_to_datadog/dd"

	"github.com/prometheus/prometheus/promql/parser"
)

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
	return id
}

func (q *Query) Build() (string, error) {
	var err error
	query := dd.Query{}

	if q.QueryType != "Azure Monitor" {
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

	query.GroupBys = q.groups()

	query.Function = q.function()

	query.Filters, err = q.filter()
	if err != nil {
		return "", err
	}

	return query.Build()
}

func (q *Query) parseExpr() error {
	var expr parser.Expr
	var err error

	expr, err = parser.ParseExpr(q.Expr)
	if err != nil {
		return err
	}

	if expr.Type() != parser.ValueTypeVector {
		return fmt.Errorf("Expression type %s note supported")
	}

	agg, ok := expr.(*parser.AggregateExpr)
	if ok {
		fmt.Println("Agg: ", agg.Op)
		expr = agg.Expr
	}
}
