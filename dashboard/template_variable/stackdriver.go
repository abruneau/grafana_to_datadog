package templatevariable

import (
	"grafana_to_datadog/grafana"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

func stackdriverConverter(tv grafana.TemplateVariable) *datadogV1.DashboardTemplateVariable {
	if tv.Name == "alignmentPeriod" {
		return nil
	}
	return defaultConverter(tv)
}
