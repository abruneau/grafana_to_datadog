package widgets

import (
	"grafana_to_datadog/dashboard/widgets/converter"
	"grafana_to_datadog/grafana"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

var displayMap = map[string]datadogV1.WidgetDisplayType{
	"line":   datadogV1.WIDGETDISPLAYTYPE_LINE,
	"bar":    datadogV1.WIDGETDISPLAYTYPE_BARS,
	"points": datadogV1.WIDGETDISPLAYTYPE_LINE,
}

func newTimeseriesDefinition(source string, panel grafana.Panel) (datadogV1.WidgetDefinition, error) {
	request, err := newTimeseriesRequest(source, panel)
	if err != nil {
		return datadogV1.WidgetDefinition{}, err
	}
	tsDefinition := datadogV1.NewTimeseriesWidgetDefinition(request, datadogV1.TIMESERIESWIDGETDEFINITIONTYPE_TIMESERIES)
	tsDefinition.SetTitle(panel.Title)
	tsDefinition.SetTitleSize("16")

	return datadogV1.TimeseriesWidgetDefinitionAsWidgetDefinition(tsDefinition), nil
}

func newTimeseriesRequest(source string, panel grafana.Panel) ([]datadogV1.TimeseriesWidgetRequest, error) {
	var widgetRequest *datadogV1.TimeseriesWidgetRequest
	var err error

	if source == "" {
		source = panel.Datasource.Type
	}

	con, err := converter.NewConverter(source)
	if err != nil {
		return nil, err
	}

	widgetRequest, err = con.NewTimeseriesWidgetRequest(panel)

	if err != nil {
		return nil, err
	}

	if panel.Type == "barchart" {
		widgetRequest.SetDisplayType(datadogV1.WIDGETDISPLAYTYPE_BARS)
	} else if panel.FieldConfig.Defaults.Custom.DrawStyle != "" {
		widgetRequest.SetDisplayType(displayMap[panel.FieldConfig.Defaults.Custom.DrawStyle])
	}
	widgetRequest.SetResponseFormat(datadogV1.FORMULAANDFUNCTIONRESPONSEFORMAT_TIMESERIES)

	return []datadogV1.TimeseriesWidgetRequest{*widgetRequest}, nil
}
