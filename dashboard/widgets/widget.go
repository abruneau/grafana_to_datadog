package widgets

import (
	"fmt"
	"grafana_to_datadog/grafana"
	"math"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

func ConvertWidget(source string, panel grafana.Panel) (datadogV1.Widget, error) {

	widget := datadogV1.NewWidgetWithDefaults()
	var definition datadogV1.WidgetDefinition
	var err error = nil

	switch panel.Type {
	case "piechart":
		definition, err = newPiechartDefinition(source, panel)
	case "row":
		definition, err = newGroupDefinition(panel)
	case "stat", "singlestat":
		definition, err = newQueryValueDefinition(source, panel)
	case "text":
		definition, err = newTextDefinition(panel)
	case "timeseries", "graph", "barchart":
		definition, err = newTimeseriesDefinition(source, panel)
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
	heigh := int64(math.Floor(float64(panel.GridPos.H) / 2))
	if heigh == 0 {
		heigh = 1
	}
	width := int64(math.Floor(float64(panel.GridPos.W) / 2))
	x := int64(math.Floor(float64(panel.GridPos.X) / 2))
	y := int64(math.Floor(float64(panel.GridPos.Y) / 2))
	return *datadogV1.NewWidgetLayout(heigh, width, x, y)
}
