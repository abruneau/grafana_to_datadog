package cloudwatch

import "github.com/mitchellh/mapstructure"

type Target struct {
	Alias       string `json:"alias"`
	Application struct {
		Filter string `json:"filter"`
	} `json:"application"`
	Datasource struct {
		Type string `json:"type"`
		UID  string `json:"uid"`
	} `json:"datasource"`
	Dimensions map[string]string `json:"dimensions"`
	Expression string            `json:"expression"`
	Functions  []interface{}     `json:"functions"`
	Group      struct {
		Filter string `json:"filter"`
	} `json:"group"`
	Host struct {
		Filter string `json:"filter"`
	} `json:"host"`
	ID   string `json:"id"`
	Item struct {
		Filter string `json:"filter"`
	} `json:"item"`
	MatchExact       bool   `json:"matchExact"`
	MetricEditorMode int    `json:"metricEditorMode"`
	MetricName       string `json:"metricName"`
	MetricQueryType  int    `json:"metricQueryType"`
	MetricQuery      struct {
		AliasBy            string        `json:"aliasBy"`
		AlignmentPeriod    string        `json:"alignmentPeriod"`
		CrossSeriesReducer string        `json:"crossSeriesReducer"`
		EditorMode         string        `json:"editorMode"`
		Filters            []string      `json:"filters"`
		GroupBys           []interface{} `json:"groupBys"`
		MetricKind         string        `json:"metricKind"`
		MetricType         string        `json:"metricType"`
		PerSeriesAligner   string        `json:"perSeriesAligner"`
		Preprocessor       string        `json:"preprocessor"`
		ProjectName        string        `json:"projectName"`
		Query              string        `json:"query"`
		ValueType          string        `json:"valueType"`
	} `json:"metricQuery"`
	Mode      int    `json:"mode"`
	Namespace string `json:"namespace"`
	Options   struct {
		ShowDisabledItems bool `json:"showDisabledItems"`
	} `json:"options"`
	Period        string   `json:"period"`
	QueryMode     string   `json:"queryMode"`
	RefID         string   `json:"refId"`
	Region        string   `json:"region"`
	SQLExpression string   `json:"sqlExpression"`
	Statistic     string   `json:"statistic"`
	Statistics    []string `json:"statistics"`
	Hide          bool     `json:"hide"`
}

func NewTarget(target map[string]interface{}) *Target {
	var output Target
	mapstructure.Decode(target, &output)
	return &output
}
