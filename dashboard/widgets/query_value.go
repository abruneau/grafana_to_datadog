package widgets

import (
	"fmt"
	"grafana_to_datadog/dashboard/widgets/cloudwatch"
	"grafana_to_datadog/dashboard/widgets/stackdriver"
	"grafana_to_datadog/grafana"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	log "github.com/sirupsen/logrus"
)

func newQueryValueDefinition(source string, panel grafana.Panel, logger *log.Entry) (datadogV1.WidgetDefinition, error) {
	request, err := newQueryValueRequest(source, panel, logger)
	if err != nil {
		return datadogV1.WidgetDefinition{}, err
	}
	qvDefinition := datadogV1.NewQueryValueWidgetDefinition(request, datadogV1.QUERYVALUEWIDGETDEFINITIONTYPE_QUERY_VALUE)
	qvDefinition.SetTitle(panel.Title)
	qvDefinition.SetTitleSize("16")

	if panel.Options.GraphMode == "area" {
		qvDefinition.TimeseriesBackground = datadogV1.NewTimeseriesBackground(datadogV1.TIMESERIESBACKGROUNDTYPE_AREA)
	}

	return datadogV1.QueryValueWidgetDefinitionAsWidgetDefinition(qvDefinition), nil
}

func newQueryValueRequest(source string, panel grafana.Panel, logger *log.Entry) ([]datadogV1.QueryValueWidgetRequest, error) {
	var widgetRequest *datadogV1.QueryValueWidgetRequest
	var err error

	if source == "" {
		source = panel.Datasource.Type
	}

	switch source {
	case "cloudwatch":
		widgetRequest, err = cloudwatch.NewQueryValueWidgetRequest(panel, logger)
	case "stackdriver":
		widgetRequest, err = stackdriver.NewQueryValueWidgetRequest(panel, logger)
	default:
		err = fmt.Errorf("unknown datasource %s", panel.Datasource.Type)
	}

	if err != nil {
		return nil, err
	}

	widgetRequest.SetResponseFormat(datadogV1.FORMULAANDFUNCTIONRESPONSEFORMAT_SCALAR)

	return []datadogV1.QueryValueWidgetRequest{*widgetRequest}, nil
}
