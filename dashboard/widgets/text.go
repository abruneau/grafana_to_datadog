package widgets

import (
	"grafana_to_datadog/grafana"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	md "github.com/JohannesKaufmann/html-to-markdown"
)

func newTextDefinition(panel grafana.Panel) (datadogV1.WidgetDefinition, error) {
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(panel.Options.Content)
	if err != nil {
		markdown = panel.Options.Content
	}
	def := datadogV1.NewNoteWidgetDefinition(markdown, datadogV1.NOTEWIDGETDEFINITIONTYPE_NOTE)
	def.SetShowTick(false)
	def.SetBackgroundColor("white")
	return datadogV1.NoteWidgetDefinitionAsWidgetDefinition(def), nil
}
