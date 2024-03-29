package stackdriver

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
	TimeSeriesList struct {
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
	} `json:"timeSeriesList"`
	QueryType string `json:"queryType"`
	RefID     string `json:"refId"`
}
