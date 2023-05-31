package widgets

import (
	"fmt"
	"grafana_to_datadog/dashboard/widgets/cloudwatch"
	"grafana_to_datadog/dashboard/widgets/stackdriver"
	"grafana_to_datadog/grafana"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	log "github.com/sirupsen/logrus"
)

func newPiechartDefinition(source string, panel grafana.Panel, logger *log.Entry) (datadogV1.WidgetDefinition, error) {
	request, err := newPiechartRequest(source, panel, logger)
	if err != nil {
		return datadogV1.WidgetDefinition{}, err
	}
	def := datadogV1.NewSunburstWidgetDefinition(request, datadogV1.SUNBURSTWIDGETDEFINITIONTYPE_SUNBURST)
	def.SetTitle(panel.Title)
	def.SetTitleSize("16")

	return datadogV1.SunburstWidgetDefinitionAsWidgetDefinition(def), nil
}

func newPiechartRequest(source string, panel grafana.Panel, logger *log.Entry) ([]datadogV1.SunburstWidgetRequest, error) {
	var widgetRequest *datadogV1.SunburstWidgetRequest
	var err error

	if source == "" {
		source = panel.Datasource.Type
	}

	switch source {
	case "cloudwatch":
		widgetRequest, err = cloudwatch.NewSunburstWidgetRequest(panel, logger)
	case "stackdriver":
		widgetRequest, err = stackdriver.NewSunburstWidgetRequest(panel, logger)
	default:
		err = fmt.Errorf("unknown datasource %s", panel.Datasource.Type)
	}

	if err != nil {
		return nil, err
	}

	widgetRequest.SetResponseFormat(datadogV1.FORMULAANDFUNCTIONRESPONSEFORMAT_SCALAR)

	return []datadogV1.SunburstWidgetRequest{*widgetRequest}, nil
}
