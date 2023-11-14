package dd

type FormulaAndFunctionMetricTransformation string

// List of FormulaAndFunctionMetricTransformation.
const (
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_ABS        FormulaAndFunctionMetricTransformation = "abs"
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_CLAMP_MIN  FormulaAndFunctionMetricTransformation = "clamp_min"
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_CLAMP_MAX  FormulaAndFunctionMetricTransformation = "clamp_max"
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_DERIV      FormulaAndFunctionMetricTransformation = "deriv"
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_LOG2       FormulaAndFunctionMetricTransformation = "log2"
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_LOG10      FormulaAndFunctionMetricTransformation = "log10"
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_DELTA      FormulaAndFunctionMetricTransformation = "delta"
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_RATE       FormulaAndFunctionMetricTransformation = "rate"
	FORMULAANDFUNCTIONMETRICTRANSFORMATION_PER_SECOND FormulaAndFunctionMetricTransformation = "per_second"
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
