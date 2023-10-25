package prometheus

type Target struct {
	Expr           string `json:"expr"`
	Format         string `json:"format"`
	IntervalFactor int    `json:"intervalFactor"`
	LegendFormat   string `json:"legendFormat"`
	RefID          string `json:"refId"`
	Step           int    `json:"step,omitempty"`
	Hide           bool   `json:"hide"`
}
