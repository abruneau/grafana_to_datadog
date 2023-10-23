package azure

import (
	"fmt"
	"grafana_to_datadog/dashboard/widgets/shared"
	"grafana_to_datadog/dd"
	"grafana_to_datadog/utilities"
	"strings"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/gobeam/stringy"
)

var (
	root          = "azure"
	agregationMap = map[string]datadogV1.FormulaAndFunctionMetricAggregation{
		"Average": "avg",
		"Sum":     "sum",
		"Maximum": "max",
		"Minimum": "min",
		"Total":   "sum",
		"Count":   "sum",
		"None":    "",
	}
	metricMap = map[string]string{
		"Compute/virtualMachines": "vm",
	}
)

type Query struct {
	*Target
	groupBy bool
}

func NewQuery(target map[string]interface{}, groupBy bool) shared.Query {
	t := shared.NewTarget[Target](target)
	query := &Query{
		Target:  t,
		groupBy: groupBy,
	}
	return query
}

func (q *Query) Id() string {
	id := q.RefID
	// if q.ID != "" {
	// 	id = q.ID
	// }
	return id
}

func (q *Query) Build() (string, error) {
	var err error
	query := dd.Query{}

	if q.QueryType != "Azure Monitor" {
		return "", fmt.Errorf("unsupported query mode %s", q.QueryType)
	}

	query.Metric, err = q.metric()
	if err != nil {
		return "", err
	}

	query.Aggregator, err = q.Aggregator()
	if err != nil {
		return "", err
	}

	query.GroupBys = q.groups()

	query.Function = q.function()

	query.Filters, err = q.filter()
	if err != nil {
		return "", err
	}

	return query.Build()
}

func (q *Query) Aggregator() (datadogV1.FormulaAndFunctionMetricAggregation, error) {
	agg, ok := agregationMap[q.AzureMonitor.Aggregation]
	if !ok {
		return "", shared.AggregationTypeError(q.AzureMonitor.Aggregation)
	}
	return agg, nil
}

func (q *Query) Formula() string {
	if q.Hide {
		return ""
	}
	return q.Id()
}

func (q *Query) metric() (string, error) {

	ns := q.AzureMonitor.MetricDefinition
	if q.AzureMonitor.MetricDefinition == "" {
		ns = q.AzureMonitor.MetricNamespace
	}

	if ns[0] == byte('$') {
		return "", fmt.Errorf("namespace as template variable not supported")
	}

	defParts := strings.Split(ns, ".")
	if len(defParts) < 2 {
		return "", fmt.Errorf("no namespace found")
	}

	var resourceType string

	resourceType, ok := metricMap[defParts[1]]
	if !ok {
		resourceType = strings.ToLower(strings.ReplaceAll(defParts[1], "/", "_"))
	}

	metricName := stringy.New(strings.ReplaceAll(q.AzureMonitor.MetricName, "/", " "))
	metric := metricName.SnakeCase().ToLower()

	return strings.Join([]string{root, resourceType, metric}, "."), nil
}

func (q *Query) groups() []string {
	variables := []string{}
	if !q.groupBy {
		return variables
	}
	for _, v := range q.AzureMonitor.Dimensions {
		if v.Text != "None" {
			dimName := stringy.New(v.Value)
			variables = append(variables, dimName.ToLower())
		}
	}
	return variables
}

// TODO: It should not be static
func (q *Query) function() dd.FormulaAndFunctionMetricFunction {
	return "as_count()"
}

func (q *Query) filter() ([]string, error) {
	filters := []string{}
	if q.AzureMonitor.ResourceGroup != "" {
		filters = append(filters, q.AzureMonitor.ResourceGroup)
	}
	if q.AzureMonitor.ResourceName != "" {
		filters = append(filters, q.AzureMonitor.ResourceName)
	}
	for _, f := range q.AzureMonitor.DimensionFilters {
		if f.Filters[0] == "*" {
			continue
		}
		label := strings.ToLower(f.Dimension)
		if len(f.Filters) == 1 {
			if f.Operator == "eq" {
				filters = append(filters, fmt.Sprintf("%s:%s", label, f.Filters[0]))
			} else if f.Operator == "ne" {
				filters = append(filters, fmt.Sprintf("!%s:%s", label, f.Filters[0]))
			} else if f.Operator == "sw" {
				filters = append(filters, fmt.Sprintf("%s:%s*", label, f.Filters[0]))
			}
		} else {
			value := strings.Join(f.Filters, ", ")
			if f.Operator == "eq" {
				filters = append(filters, fmt.Sprintf("%s IN (%s)", label, value))
			} else if f.Operator == "ne" {
				filters = append(filters, fmt.Sprintf("%s NOT IN (%s)", label, value))
			} else if f.Operator == "sw" {
				fs := utilities.Map(f.Filters, func(item string) string { return fmt.Sprintf("%s*", item) })
				value = strings.Join(fs, ", ")
				filters = append(filters, fmt.Sprintf("%s IN (%s)", label, value))
			}
		}
	}
	return filters, nil
}
