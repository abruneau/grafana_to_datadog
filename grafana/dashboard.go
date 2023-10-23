package grafana

type Dashboard struct {
	Inputs   []interface{} `json:"__inputs"`
	Elements []interface{} `json:"__elements"`
	Requires []struct {
		Type    string `json:"type"`
		ID      string `json:"id"`
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"__requires"`
	Annotations struct {
		List []struct {
			BuiltIn    int    `json:"builtIn"`
			Datasource string `json:"datasource"`
			Enable     bool   `json:"enable"`
			Hide       bool   `json:"hide"`
			IconColor  string `json:"iconColor"`
			Name       string `json:"name"`
			Target     struct {
				Limit    int           `json:"limit"`
				MatchAny bool          `json:"matchAny"`
				Tags     []interface{} `json:"tags"`
				Type     string        `json:"type"`
			} `json:"target"`
			Type string `json:"type"`
		} `json:"list"`
	} `json:"annotations"`
	Description          string        `json:"description"`
	Editable             bool          `json:"editable"`
	FiscalYearStartMonth int           `json:"fiscalYearStartMonth"`
	GnetID               int           `json:"gnetId"`
	GraphTooltip         int           `json:"graphTooltip"`
	ID                   interface{}   `json:"id"`
	Iteration            int64         `json:"iteration"`
	Links                []interface{} `json:"links"`
	LiveNow              bool          `json:"liveNow"`
	Panels               []Panel       `json:"panels"`
	Refresh              bool          `json:"refresh"`
	SchemaVersion        int           `json:"schemaVersion"`
	Style                string        `json:"style"`
	Tags                 []string      `json:"tags"`
	Templating           struct {
		List []TemplateVariable `json:"list"`
	} `json:"templating"`
	Time struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"time"`
	Timepicker struct {
		RefreshIntervals []string `json:"refresh_intervals"`
		TimeOptions      []string `json:"time_options"`
	} `json:"timepicker"`
	Timezone  string `json:"timezone"`
	Title     string `json:"title"`
	UID       string `json:"uid"`
	Version   int    `json:"version"`
	WeekStart string `json:"weekStart"`
}

type TemplateVariable struct {
	Current struct {
	} `json:"current"`
	Hide        int           `json:"hide"`
	IncludeAll  bool          `json:"includeAll"`
	Label       string        `json:"label"`
	Multi       bool          `json:"multi"`
	Name        string        `json:"name"`
	Options     []interface{} `json:"options"`
	Query       string        `json:"query"`
	QueryValue  string        `json:"queryValue,omitempty"`
	Refresh     int           `json:"refresh"`
	Regex       string        `json:"regex"`
	SkipURLSync bool          `json:"skipUrlSync"`
	Type        string        `json:"type"`
	Datasource  struct {
		UID string `json:"uid"`
	} `json:"datasource,omitempty"`
	Definition     string `json:"definition,omitempty"`
	Sort           int    `json:"sort,omitempty"`
	TagValuesQuery string `json:"tagValuesQuery,omitempty"`
	TagsQuery      string `json:"tagsQuery,omitempty"`
	UseTags        bool   `json:"useTags,omitempty"`
}
