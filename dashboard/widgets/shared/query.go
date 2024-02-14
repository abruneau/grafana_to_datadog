package shared

import (
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/mitchellh/mapstructure"
)

func NewTarget[T interface{}](target map[string]interface{}) *T {
	var output T
	mapstructure.Decode(target, &output)
	return &output
}

type Request struct {
	Formulas []string
	Queries  []struct {
		Name        string
		Query       string
		Aggregation datadogV1.FormulaAndFunctionMetricAggregation
	}
}
