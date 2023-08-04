package grafana

type Panel struct {
	Datasource struct {
		Type string `json:"type"`
		UID  string `json:"uid"`
	} `json:"datasource,omitempty"`
	FieldConfig struct {
		Defaults struct {
			Color struct {
				Mode string `json:"mode"`
			} `json:"color"`
			Custom struct {
				AxisLabel     string `json:"axisLabel"`
				AxisPlacement string `json:"axisPlacement"`
				BarAlignment  int    `json:"barAlignment"`
				DrawStyle     string `json:"drawStyle"`
				FillOpacity   int    `json:"fillOpacity"`
				GradientMode  string `json:"gradientMode"`
				HideFrom      struct {
					Legend  bool `json:"legend"`
					Tooltip bool `json:"tooltip"`
					Viz     bool `json:"viz"`
				} `json:"hideFrom"`
				LineInterpolation string `json:"lineInterpolation"`
				LineWidth         int    `json:"lineWidth"`
				PointSize         int    `json:"pointSize"`
				ScaleDistribution struct {
					Type string `json:"type"`
				} `json:"scaleDistribution"`
				ShowPoints string `json:"showPoints"`
				SpanNulls  bool   `json:"spanNulls"`
				Stacking   struct {
					Group string `json:"group"`
					Mode  string `json:"mode"`
				} `json:"stacking"`
				ThresholdsStyle struct {
					Mode string `json:"mode"`
				} `json:"thresholdsStyle"`
			} `json:"custom"`
			Mappings   []interface{} `json:"mappings"`
			Min        int           `json:"min"`
			Thresholds struct {
				Mode  string `json:"mode"`
				Steps []struct {
					Color string      `json:"color"`
					Value interface{} `json:"value"`
				} `json:"steps"`
			} `json:"thresholds"`
			Unit string `json:"unit"`
		} `json:"defaults"`
		Overrides []struct {
			Matcher struct {
				ID      string `json:"id"`
				Options string `json:"options"`
			} `json:"matcher"`
			Properties []struct {
				ID    string `json:"id"`
				Value string `json:"value"`
			} `json:"properties"`
		} `json:"overrides"`
	} `json:"fieldConfig,omitempty"`
	GridPos struct {
		H int `json:"h"`
		W int `json:"w"`
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"gridPos"`
	ID      int           `json:"id"`
	Links   []interface{} `json:"links,omitempty"`
	Options struct {
		Legend struct {
			Calcs       []string `json:"calcs"`
			DisplayMode string   `json:"displayMode"`
			Placement   string   `json:"placement"`
		} `json:"legend"`
		Tooltip struct {
			Mode string `json:"mode"`
			Sort string `json:"sort"`
		} `json:"tooltip"`
		Content   string `json:"content"`
		GraphMode string `json:"graphMode"`
		Mode      string `json:"mode"`
	} `json:"options,omitempty"`
	Panels        []Panel
	PluginVersion string                   `json:"pluginVersion"`
	Targets       []map[string]interface{} `json:"targets,omitempty"`
	Title         string                   `json:"title"`
	Type          string                   `json:"type"`
	Editable      bool                     `json:"editable,omitempty"`
	Error         bool                     `json:"error,omitempty"`
}
