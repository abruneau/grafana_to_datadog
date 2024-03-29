package converter

import (
	"fmt"
	"grafana_to_datadog/dashboard/widgets/converter/azure"
	"grafana_to_datadog/dashboard/widgets/converter/cloudwatch"
	"grafana_to_datadog/dashboard/widgets/converter/prometheus"
	"grafana_to_datadog/dashboard/widgets/converter/stackdriver"
	"grafana_to_datadog/dashboard/widgets/shared"
	"grafana_to_datadog/grafana"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

type NewQueryFunction func(target map[string]interface{}, groupBy bool) (shared.Request, error)

var sourceMapper = map[string]NewQueryFunction{
	"grafana-azure-monitor-datasource": azure.NewQuery,
	"cloudwatch":                       cloudwatch.NewQuery,
	"stackdriver":                      stackdriver.NewQuery,
	"prometheus":                       prometheus.NewQuery,
	"loki":                             prometheus.NewQuery,
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

func (c *Converter) parseTargets(panel grafana.Panel, aggregate bool, groupBy bool) (queries []datadogV1.FormulaAndFunctionQueryDefinition, formulas []datadogV1.WidgetFormula, err error) {
	queries = []datadogV1.FormulaAndFunctionQueryDefinition{}
	formulas = []datadogV1.WidgetFormula{}

	for _, t := range panel.Targets {
		r, err := c.newQuery(t, groupBy)

		if err != nil {
			return queries, formulas, err
		}

		for _, f := range r.Queries {
			q := datadogV1.NewFormulaAndFunctionMetricQueryDefinition("metrics", f.Name, f.Query)
			queries = append(queries, datadogV1.FormulaAndFunctionMetricQueryDefinitionAsFormulaAndFunctionQueryDefinition(q))
			if aggregate {
				q.SetAggregator(f.Aggregation)
			}
		}

		for _, f := range r.Formulas {
			formulas = append(formulas, *datadogV1.NewWidgetFormula(f))
		}
	}

	if len(formulas) == 0 {
		return queries, formulas, shared.NoValidFormulaError(panel.Title)
	}

	return queries, formulas, nil
}

func (c *Converter) NewTimeseriesWidgetRequest(panel grafana.Panel) (*datadogV1.TimeseriesWidgetRequest, error) {
	var err error
	widgetRequest := datadogV1.NewTimeseriesWidgetRequest()
	widgetRequest.Queries, widgetRequest.Formulas, err = c.parseTargets(panel, false, true)

	return widgetRequest, err
}

func (c *Converter) NewQueryValueWidgetRequest(panel grafana.Panel) (*datadogV1.QueryValueWidgetRequest, error) {
	var err error
	widgetRequest := datadogV1.NewQueryValueWidgetRequest()
	widgetRequest.Queries, widgetRequest.Formulas, err = c.parseTargets(panel, true, false)

	return widgetRequest, err
}

func (c *Converter) NewSunburstWidgetRequest(panel grafana.Panel) (*datadogV1.SunburstWidgetRequest, error) {
	var err error
	widgetRequest := datadogV1.NewSunburstWidgetRequest()
	widgetRequest.Queries, widgetRequest.Formulas, err = c.parseTargets(panel, true, true)

	return widgetRequest, err
}
