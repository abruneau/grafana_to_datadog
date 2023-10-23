package azure

type Target struct {
	Datasource struct {
		Type string `json:"type"`
		UID  string `json:"uid"`
	} `json:"datasource"`
	RefID        string `json:"refId"`
	QueryType    string `json:"queryType"`
	AzureMonitor struct {
		Aggregation         string `json:"aggregation"`
		TimeGrain           string `json:"timeGrain"`
		AllowedTimeGrainsMs []int  `json:"allowedTimeGrainsMs"`
		MetricNamespace     string `json:"metricNamespace"`
		MetricDefinition    string `json:"metricDefinition"`
		Region              string `json:"region"`
		Resources           []struct {
			Subscription    string `json:"subscription"`
			ResourceGroup   string `json:"resourceGroup"`
			MetricNamespace string `json:"metricNamespace"`
			ResourceName    string `json:"resourceName"`
			Region          string `json:"region"`
		} `json:"resources"`
		MetricName       string `json:"metricName"`
		DimensionFilter  string `json:"dimentionFilter"`
		DimensionFilters []struct {
			Dimension string   `json:"dimension"`
			Operator  string   `json:"operator"`
			Filters   []string `json:"filters"`
		} `json:"dimensionFilters"`
		Dimensions []struct {
			Text  string `json:"text"`
			Value string `json:"value"`
		} `json:"dimensions"`
		ResourceGroup string `json:"resourceGroup"`
		ResourceName  string `json:"resourceName"`
	} `json:"azureMonitor"`
	Subscription string `json:"subscription"`
	Hide         bool   `json:"hide"`
}
