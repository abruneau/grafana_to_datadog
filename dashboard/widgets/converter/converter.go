package converter

import (
	"fmt"
	"grafana_to_datadog/dashboard/widgets/converter/cloudwatch"
	"grafana_to_datadog/dashboard/widgets/converter/stackdriver"
	"grafana_to_datadog/dashboard/widgets/shared"
	"grafana_to_datadog/grafana"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

type NewQueryFunction func(target map[string]interface{}) shared.Query

var sourceMapper = map[string]NewQueryFunction{
	"cloudwatch":  cloudwatch.NewQuery,
	"stackdriver": stackdriver.NewQuery,
}

type Converter struct {
	source   string
	newQuery NewQueryFunction
}

func NewConverter(source string) (*Converter, error) {
	var ok bool
	conv := &Converter{source: source}
	conv.newQuery, ok = sourceMapper[source]
	if !ok {
		return nil, fmt.Errorf("unknown datasource %s", source)
	}

	return conv, nil
}

func (c *Converter) parseTargets(panel grafana.Panel, aggregate bool) (queries []datadogV1.FormulaAndFunctionQueryDefinition, formulas []datadogV1.WidgetFormula, err error) {
	queries = []datadogV1.FormulaAndFunctionQueryDefinition{}
	formulas = []datadogV1.WidgetFormula{}

	for _, t := range panel.Targets {
		query := c.newQuery(t)

		id := query.Id()
		targetQuery, err := query.Build()
		if err != nil {
			return queries, formulas, err
		}

		q := datadogV1.NewFormulaAndFunctionMetricQueryDefinition("metrics", id, targetQuery)
		if aggregate {
			agg, _ := query.Aggregator()
			q.SetAggregator(agg)
		}
		queries = append(queries, datadogV1.FormulaAndFunctionMetricQueryDefinitionAsFormulaAndFunctionQueryDefinition(q))
		formulas = append(formulas, *datadogV1.NewWidgetFormula(query.Formula()))
	}

	if len(formulas) == 0 {
		return queries, formulas, shared.NoValidFormulaError(panel.Title)
	}

	return queries, formulas, nil
}

func (c *Converter) NewTimeseriesWidgetRequest(panel grafana.Panel) (*datadogV1.TimeseriesWidgetRequest, error) {
	var err error
	widgetRequest := datadogV1.NewTimeseriesWidgetRequest()
	widgetRequest.Queries, widgetRequest.Formulas, err = c.parseTargets(panel, false)

	return widgetRequest, err
}

func (c *Converter) NewQueryValueWidgetRequest(panel grafana.Panel) (*datadogV1.QueryValueWidgetRequest, error) {
	var err error
	widgetRequest := datadogV1.NewQueryValueWidgetRequest()
	widgetRequest.Queries, widgetRequest.Formulas, err = c.parseTargets(panel, true)

	return widgetRequest, err
}

func (c *Converter) NewSunburstWidgetRequest(panel grafana.Panel) (*datadogV1.SunburstWidgetRequest, error) {
	var err error
	widgetRequest := datadogV1.NewSunburstWidgetRequest()
	widgetRequest.Queries, widgetRequest.Formulas, err = c.parseTargets(panel, true)

	return widgetRequest, err
}
