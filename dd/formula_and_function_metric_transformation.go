package dd

type FormulaAndFunctionMetricTransformation string

// List of FormulaAndFunctionMetricTransformation.
const (
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_ABS          FormulaAndFunctionMetricTransformation = "abs"
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_CLAMP_MIN    FormulaAndFunctionMetricTransformation = "clamp_min"
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_CLAMP_MAX    FormulaAndFunctionMetricTransformation = "clamp_max"
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_DERIV        FormulaAndFunctionMetricTransformation = "deriv"
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_LOG2         FormulaAndFunctionMetricTransformation = "log2"
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_LOG10        FormulaAndFunctionMetricTransformation = "log10"
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_DELTA        FormulaAndFunctionMetricTransformation = "delta"
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_RATE         FormulaAndFunctionMetricTransformation = "rate"
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_PER_SECOND   FormulaAndFunctionMetricTransformation = "per_second"
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_HOUR_BEFORE  FormulaAndFunctionMetricTransformation = "hour_before"
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_DAY_BEFORE   FormulaAndFunctionMetricTransformation = "day_before"
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_WEEK_BEFORE  FormulaAndFunctionMetricTransformation = "week_before"
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_MONTH_BEFORE FormulaAndFunctionMetricTransformation = "month_before"
)

var allowedFormulaAndFunctionMetricTransformationEnumValues = []FormulaAndFunctionMetricTransformation{
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_ABS,
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_CLAMP_MIN,
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_CLAMP_MAX,
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_DERIV,
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_LOG2,
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_LOG10,
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_DELTA,
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_RATE,
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_PER_SECOND,
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_HOUR_BEFORE,
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_DAY_BEFORE,
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_WEEK_BEFORE,
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_MONTH_BEFORE,
}

// GetAllowedValues reeturns the list of possible values.
func (v *FormulaAndFunctionMetricTransformation) GetAllowedValues() []FormulaAndFunctionMetricTransformation {
	return allowedFormulaAndFunctionMetricTransformationEnumValues
}

// IsValid return true if the value is valid for the enum, false otherwise.
func (v FormulaAndFunctionMetricTransformation) IsValid() bool {
	for _, existing := range allowedFormulaAndFunctionMetricTransformationEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}
