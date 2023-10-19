package shared

import (
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/mitchellh/mapstructure"
)

type Query interface {
	Id() string
	Build() (string, error)
	Aggregator() (datadogV1.FormulaAndFunctionMetricAggregation, error)
	Formula() string
}

func NewTarget[T interface{}](target map[string]interface{}) *T {
	var output T
	mapstructure.Decode(target, &output)
	return &output
}
