package dashboard

import (
	"fmt"
	"grafana_to_datadog/dashboard/widgets"
	"grafana_to_datadog/grafana"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

var templateVariablesBlacklist = []string{"alignmentPeriod"}

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
		if v.Type == "query" && !slices.Contains(templateVariablesBlacklist, v.Name) {
			tv := datadogV1.NewDashboardTemplateVariable(v.Name)
			tv.SetPrefix(v.Name)
			c.templateVariables = append(c.templateVariables, *tv)
		} else if v.Type == "datasource" {
			c.datasource = v.Query
		}
	}
}

func (c *dashboardConvertor) processPanel(panel grafana.Panel) *datadogV1.Widget {
	widget, err := widgets.ConvertWidget(c.datasource, panel)
	if err != nil {
		c.logger.Error(err)
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
