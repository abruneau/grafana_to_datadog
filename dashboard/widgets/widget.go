package widgets

import (
	"fmt"
	"grafana_to_datadog/grafana"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	log "github.com/sirupsen/logrus"
)

func ConvertWidget(source string, panel grafana.Panel, logger *log.Entry) (datadogV1.Widget, error) {

	widget := datadogV1.NewWidgetWithDefaults()
	var definition datadogV1.WidgetDefinition
	var err error = nil

	switch panel.Type {
	case "timeseries", "graph", "barchart":
		definition, err = newTimeseriesDefinition(source, panel, logger)
	case "text":
		definition, err = newTextDefinition(panel)
	case "stat":
		definition, err = newQueryValueDefinition(source, panel, logger)
	case "row":
		definition, err = newGroupDefinition(panel)
	default:
		err = fmt.Errorf("unkown type: %s", panel.Type)
	}

	if err != nil {
		return *widget, err
	}

	widget.SetDefinition(definition)
	widget.SetLayout(convertGrid(panel))

	return *widget, err
}

func convertGrid(panel grafana.Panel) datadogV1.WidgetLayout {
	return *datadogV1.NewWidgetLayout(convertMesures(panel.GridPos.H), convertMesures(panel.GridPos.W), convertMesures(panel.GridPos.X), convertMesures(panel.GridPos.Y))
}

func convertMesures(value int) int64 {
	if value == 0 {
		return int64(value)
	}
	if value%2 == 0 {
		return int64(value / 2)
	} else {
		return int64((value + 1) / 2)
	}
}
