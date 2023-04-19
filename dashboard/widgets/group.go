package widgets

import (
	"grafana_to_datadog/grafana"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

func newGroupDefinition(panel grafana.Panel) (datadogV1.WidgetDefinition, error) {
	def := datadogV1.NewGroupWidgetDefinition(datadogV1.WIDGETLAYOUTTYPE_ORDERED, datadogV1.GROUPWIDGETDEFINITIONTYPE_GROUP, []datadogV1.Widget{})
	def.SetTitle(panel.Title)
	def.SetShowTitle(true)

	return datadogV1.GroupWidgetDefinitionAsWidgetDefinition(def), nil
}
