package stackdriver

import (
	"github.com/mitchellh/mapstructure"
)

type Target struct {
	Datasource struct {
		UID string `json:"uid"`
	} `json:"datasource"`
	Hide        bool `json:"hide"`
	MetricQuery struct {
		AliasBy          string   `json:"aliasBy"`
		AlignmentPeriod  string   `json:"alignmentPeriod"`
		Filters          []string `json:"filters"`
		GroupBys         []string `json:"groupBys"`
		MetricKind       string   `json:"metricKind"`
		MetricType       string   `json:"metricType"`
		PerSeriesAligner string   `json:"perSeriesAligner"`
		ProjectName      string   `json:"projectName"`
		Unit             string   `json:"unit"`
		ValueType        string   `json:"valueType"`
	} `json:"metricQuery"`
	QueryType string `json:"queryType"`
	RefID     string `json:"refId"`
}

func NewTarget(target map[string]interface{}) *Target {
	var output Target
	mapstructure.Decode(target, &output)
	return &output
}
