package widgets

import (
	"fmt"
	"grafana_to_datadog/grafana"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	md "github.com/JohannesKaufmann/html-to-markdown"
)

func newTextDefinition(panel grafana.Panel) (datadogV1.WidgetDefinition, error) {
	var markdown string
	if panel.Title != "" {
		markdown = fmt.Sprintf("# %s\n", panel.Title)
	}
	if panel.Options.Mode == "markdown" {
		markdown = markdown + panel.Options.Content
	} else if panel.Options.Mode == "html" {
		converter := md.NewConverter("", true, nil)
		m, err := converter.ConvertString(panel.Options.Content)
		if err != nil {
			markdown = markdown + panel.Options.Content
		}
		markdown = markdown + m

	}
	markdown = markdown + panel.Options.Content
	def := datadogV1.NewNoteWidgetDefinition(markdown, datadogV1.NOTEWIDGETDEFINITIONTYPE_NOTE)
	def.SetShowTick(false)
	def.SetBackgroundColor("white")
	return datadogV1.NoteWidgetDefinitionAsWidgetDefinition(def), nil
}
