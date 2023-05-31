package stackdriver

import (
	"fmt"
	"grafana_to_datadog/grafana"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	log "github.com/sirupsen/logrus"
)

func NewTimeseriesWidgetRequest(panel grafana.Panel, logger *log.Entry) (*datadogV1.TimeseriesWidgetRequest, error) {
	widgetRequest := datadogV1.NewTimeseriesWidgetRequest()

	for _, t := range panel.Targets {
		target := NewTarget(t)
		query := Query{
			target,
		}

		id := query.id()
		targetQuery, err := query.build()
		if err != nil {
			return nil, err
		}

		q := datadogV1.NewFormulaAndFunctionMetricQueryDefinition("metrics", id, targetQuery)
		widgetRequest.Queries = append(widgetRequest.Queries, datadogV1.FormulaAndFunctionMetricQueryDefinitionAsFormulaAndFunctionQueryDefinition(q))

		widgetRequest.Formulas = append(widgetRequest.Formulas, *datadogV1.NewWidgetFormula(query.formula()))
	}

	if len(widgetRequest.Formulas) == 0 {
		return nil, fmt.Errorf("%s no valid formula found", panel.Title)
	}

	return widgetRequest, nil
}

func NewQueryValueWidgetRequest(panel grafana.Panel, logger *log.Entry) (*datadogV1.QueryValueWidgetRequest, error) {
	widgetRequest := datadogV1.NewQueryValueWidgetRequest()

	for _, t := range panel.Targets {
		target := NewTarget(t)
		query := Query{
			target,
		}
		id := query.id()
		targetQuery, err := query.build()
		if err != nil {
			return nil, err
		}
		q := datadogV1.NewFormulaAndFunctionMetricQueryDefinition("metrics", id, targetQuery)
		agg, _ := query.aggregator()
		q.SetAggregator(agg)
		widgetRequest.Queries = append(widgetRequest.Queries, datadogV1.FormulaAndFunctionMetricQueryDefinitionAsFormulaAndFunctionQueryDefinition(q))
		widgetRequest.Formulas = append(widgetRequest.Formulas, *datadogV1.NewWidgetFormula(query.formula()))

	}

	if len(widgetRequest.Formulas) == 0 {
		return nil, fmt.Errorf("%s no valid formula found", panel.Title)
	}

	return widgetRequest, nil
}

func NewSunburstWidgetRequest(panel grafana.Panel, logger *log.Entry) (*datadogV1.SunburstWidgetRequest, error) {
	widgetRequest := datadogV1.NewSunburstWidgetRequest()

	for _, t := range panel.Targets {
		target := NewTarget(t)
		query := Query{
			target,
		}
		id := query.id()
		targetQuery, err := query.build()
		if err != nil {
			return nil, err
		}
		q := datadogV1.NewFormulaAndFunctionMetricQueryDefinition("metrics", id, targetQuery)
		agg, _ := query.aggregator()
		q.SetAggregator(agg)
		widgetRequest.Queries = append(widgetRequest.Queries, datadogV1.FormulaAndFunctionMetricQueryDefinitionAsFormulaAndFunctionQueryDefinition(q))
		widgetRequest.Formulas = append(widgetRequest.Formulas, *datadogV1.NewWidgetFormula(query.formula()))

	}

	if len(widgetRequest.Formulas) == 0 {
		return nil, fmt.Errorf("%s no valid formula found", panel.Title)
	}

	return widgetRequest, nil
}
