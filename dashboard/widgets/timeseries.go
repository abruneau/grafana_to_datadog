package widgets

import (
	"fmt"
	"grafana_to_datadog/dashboard/widgets/cloudwatch"
	"grafana_to_datadog/dashboard/widgets/stackdriver"
	"grafana_to_datadog/grafana"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	log "github.com/sirupsen/logrus"
)

var displayMap = map[string]datadogV1.WidgetDisplayType{
	"line": datadogV1.WIDGETDISPLAYTYPE_LINE,
}

func newTimeseriesDefinition(source string, panel grafana.Panel, logger *log.Entry) (datadogV1.WidgetDefinition, error) {
	request, err := newTimeseriesRequest(source, panel, logger)
	if err != nil {
		return datadogV1.WidgetDefinition{}, err
	}
	tsDefinition := datadogV1.NewTimeseriesWidgetDefinition(request, datadogV1.TIMESERIESWIDGETDEFINITIONTYPE_TIMESERIES)
	tsDefinition.SetTitle(panel.Title)
	tsDefinition.SetTitleSize("16")

	return datadogV1.TimeseriesWidgetDefinitionAsWidgetDefinition(tsDefinition), nil
}

func newTimeseriesRequest(source string, panel grafana.Panel, logger *log.Entry) ([]datadogV1.TimeseriesWidgetRequest, error) {
	var widgetRequest *datadogV1.TimeseriesWidgetRequest
	var err error

	if source == "" {
		source = panel.Datasource.Type
	}

	switch source {
	case "cloudwatch":
		widgetRequest, err = cloudwatch.NewTimeseriesWidgetRequest(panel, logger)
	case "stackdriver":
		widgetRequest, err = stackdriver.NewTimeseriesWidgetRequest(panel, logger)
	default:
		err = fmt.Errorf("unknown datasource %s", panel.Datasource.Type)
	}

	if err != nil {
		return nil, err
	}

	if panel.FieldConfig.Defaults.Custom.DrawStyle != "" {
		widgetRequest.SetDisplayType(displayMap[panel.FieldConfig.Defaults.Custom.DrawStyle])
	}
	widgetRequest.SetResponseFormat(datadogV1.FORMULAANDFUNCTIONRESPONSEFORMAT_TIMESERIES)

	return []datadogV1.TimeseriesWidgetRequest{*widgetRequest}, nil
}
