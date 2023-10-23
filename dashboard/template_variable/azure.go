package templatevariable

import (
	"grafana_to_datadog/grafana"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

func azureConverter(tv grafana.TemplateVariable) *datadogV1.DashboardTemplateVariable {
	if tv.Name == "ResourceName" {
		ddtv := datadogV1.NewDashboardTemplateVariable(tv.Name)
		ddtv.SetPrefix("name")
		return ddtv
	}
	return defaultConverter(tv)
}
