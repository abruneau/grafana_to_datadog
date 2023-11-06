package dashboard

import (
	"fmt"
	templatevariable "grafana_to_datadog/dashboard/template_variable"
	"grafana_to_datadog/dashboard/widgets"
	"grafana_to_datadog/grafana"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	log "github.com/sirupsen/logrus"
)

type dashboardConvertor struct {
	graf              *grafana.Dashboard
	logger            *log.Entry
	templateVariables []datadogV1.DashboardTemplateVariable
	datasource        string
	widgets           []datadogV1.Widget
}

func ConvertDashboard(graf *grafana.Dashboard, logger *log.Entry) *datadogV1.Dashboard {
	convertor := &dashboardConvertor{
		graf:   graf,
		logger: logger,
	}

	convertor.init()
	return convertor.build()
}

func (c *dashboardConvertor) init() {
	c.extractTemplateVariables()
	c.extractWidgets()
}

func (c *dashboardConvertor) build() *datadogV1.Dashboard {
	dash := datadogV1.NewDashboard(datadogV1.DASHBOARDLAYOUTTYPE_ORDERED, c.graf.Title, c.widgets)
	description := fmt.Sprintf("%s\n\ngenerated with https://github.com/abruneau/grafana_to_datadog", c.graf.Description)
	dash.Description.Set(&description)
	dash.TemplateVariables = c.templateVariables
	return dash
}

func (c *dashboardConvertor) extractTemplateVariables() {
	for _, v := range c.graf.Templating.List {
		if v.Type == "datasource" {
			c.datasource = v.Query
		} else if v.Type == "query" {
			tv := templatevariable.GetTemplateVariable(c.datasource, v)
			if tv != nil {
				c.templateVariables = append(c.templateVariables, *tv)
			}
		}
	}
}

func (c *dashboardConvertor) processPanel(panel grafana.Panel) *datadogV1.Widget {
	widget, err := widgets.ConvertWidget(c.datasource, panel)
	if err != nil {
		c.logger.WithField("widget", panel.Title).Error(err)
		return nil
	}

	for _, p := range panel.Panels {
		childWidget := c.processPanel(p)
		if childWidget == nil {
			continue
		}
		widget.Definition.GroupWidgetDefinition.Widgets = append(widget.Definition.GroupWidgetDefinition.Widgets, *childWidget)
	}

	return &widget
}

func (c *dashboardConvertor) extractWidgets() {
	var lastGroupWidget *datadogV1.Widget
	for _, panel := range c.graf.Panels {
		widget := c.processPanel(panel)

		if widget == nil {
			continue
		}

		if lastGroupWidget != nil && widget.Definition.GroupWidgetDefinition == nil {
			lastGroupWidget.Definition.GroupWidgetDefinition.Widgets = append(lastGroupWidget.Definition.GroupWidgetDefinition.Widgets, *widget)
		} else {
			c.widgets = append(c.widgets, *widget)
		}

		if widget.Definition.GroupWidgetDefinition != nil {
			lastGroupWidget = widget
		}
	}
}
