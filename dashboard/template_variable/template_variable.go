package templatevariable

import (
	"grafana_to_datadog/grafana"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/gobeam/stringy"
)

var sourceMapper = map[string]func(grafana.TemplateVariable) *datadogV1.DashboardTemplateVariable{
	"stackdriver":                      stackdriverConverter,
	"grafana-azure-monitor-datasource": azureConverter,
}

func defaultConverter(tv grafana.TemplateVariable) *datadogV1.DashboardTemplateVariable {
	tvName := stringy.New(tv.Name)
	nomalizedName := tvName.SnakeCase().ToLower()
	ddtv := datadogV1.NewDashboardTemplateVariable(tv.Name)
	ddtv.SetPrefix(nomalizedName)
	return ddtv
}

func GetTemplateVariable(source string, tv grafana.TemplateVariable) *datadogV1.DashboardTemplateVariable {

	f, ok := sourceMapper[source]
	if ok {
		return f(tv)
	}

	return defaultConverter(tv)
}
