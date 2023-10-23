package widgets

import (
	"grafana_to_datadog/dashboard/widgets/converter"
	"grafana_to_datadog/grafana"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

func newQueryValueDefinition(source string, panel grafana.Panel) (datadogV1.WidgetDefinition, error) {
	request, err := newQueryValueRequest(source, panel)
	if err != nil {
		return datadogV1.WidgetDefinition{}, err
	}
	qvDefinition := datadogV1.NewQueryValueWidgetDefinition(request, datadogV1.QUERYVALUEWIDGETDEFINITIONTYPE_QUERY_VALUE)
	qvDefinition.SetTitle(panel.Title)
	qvDefinition.SetTitleSize("16")
	qvDefinition.SetAutoscale(true)
	qvDefinition.SetPrecision(2)

	if panel.Options.GraphMode == "area" {
		qvDefinition.TimeseriesBackground = datadogV1.NewTimeseriesBackground(datadogV1.TIMESERIESBACKGROUNDTYPE_AREA)
	}

	return datadogV1.QueryValueWidgetDefinitionAsWidgetDefinition(qvDefinition), nil
}

func newQueryValueRequest(source string, panel grafana.Panel) ([]datadogV1.QueryValueWidgetRequest, error) {
	var widgetRequest *datadogV1.QueryValueWidgetRequest
	var err error

	if source == "" {
		source = panel.Datasource.Type
	}

	con, err := converter.NewConverter(source)
	if err != nil {
		return nil, err
	}

	widgetRequest, err = con.NewQueryValueWidgetRequest(panel)

	if err != nil {
		return nil, err
	}

	widgetRequest.SetResponseFormat(datadogV1.FORMULAANDFUNCTIONRESPONSEFORMAT_SCALAR)

	return []datadogV1.QueryValueWidgetRequest{*widgetRequest}, nil
}
