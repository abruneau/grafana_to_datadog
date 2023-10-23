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
	var mode, content string

	mode = panel.Mode
	if panel.Mode == "" {
		mode = panel.Options.Mode
	}

	content = panel.Content
	if panel.Content == "" {
		content = panel.Options.Content
	}

	if mode == "markdown" {
		markdown = markdown + content
	} else if mode == "html" {
		converter := md.NewConverter("", true, nil)
		m, err := converter.ConvertString(content)
		if err != nil {
			markdown = markdown + content
		}
		markdown = markdown + m

	}
	// markdown = markdown + content
	def := datadogV1.NewNoteWidgetDefinition(markdown, datadogV1.NOTEWIDGETDEFINITIONTYPE_NOTE)
	def.SetShowTick(false)
	def.SetBackgroundColor("white")
	return datadogV1.NoteWidgetDefinitionAsWidgetDefinition(def), nil
}
