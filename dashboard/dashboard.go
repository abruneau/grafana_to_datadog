package dashboard

import (
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
	dash := datadogV1.NewDashboard("ordered", c.graf.Title, c.widgets)
	dash.Description.Set(&c.graf.Description)
	dash.TemplateVariables = c.templateVariables
	return dash
}

func (c *dashboardConvertor) extractTemplateVariables() {
	for _, v := range c.graf.Templating.List {
		if v.Type == "query" {
			tv := datadogV1.NewDashboardTemplateVariable(v.Name)
			tv.SetPrefix(v.Name)
			c.templateVariables = append(c.templateVariables, *tv)
		} else if v.Type == "datasource" {
			c.datasource = v.Query
		}
	}
}

func (c *dashboardConvertor) extractWidgets() {
	var lastGroupWidget *datadogV1.Widget
	for _, panel := range c.graf.Panels {
		widget, err := widgets.ConvertWidget(c.datasource, panel, c.logger)
		if err != nil {
			c.logger.Error(err)
		} else {

			if lastGroupWidget != nil && widget.Definition.GroupWidgetDefinition == nil {
				lastGroupWidget.Definition.GroupWidgetDefinition.Widgets = append(lastGroupWidget.Definition.GroupWidgetDefinition.Widgets, widget)
			} else {
				c.widgets = append(c.widgets, widget)
			}

		}
		if widget.Definition.GroupWidgetDefinition != nil {
			lastGroupWidget = &widget
		}
	}
}
