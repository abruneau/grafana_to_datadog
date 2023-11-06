package dd

// FormulaAndFunctionMetricFunction The functions methods available for metrics queries.
type FormulaAndFunctionMetricFunction string

// List of FormulaAndFunctionMetricFunction.
const (
	FORMULAANDFUNCTIONMETRICFUNCTION_COUNT FormulaAndFunctionMetricFunction = "as_count()"
	FORMULAANDFUNCTIONMETRICFUNCTION_RATE  FormulaAndFunctionMetricFunction = "as_rate()"
)

var allowedFormulaAndFunctionMetricFunctionEnumValues = []FormulaAndFunctionMetricFunction{
	FORMULAANDFUNCTIONMETRICFUNCTION_COUNT,
	FORMULAANDFUNCTIONMETRICFUNCTION_RATE,
}

// GetAllowedValues reeturns the list of possible values.
func (v *FormulaAndFunctionMetricFunction) GetAllowedValues() []FormulaAndFunctionMetricFunction {
	return allowedFormulaAndFunctionMetricFunctionEnumValues
}

// IsValid return true if the value is valid for the enum, false otherwise.
func (v FormulaAndFunctionMetricFunction) IsValid() bool {
	for _, existing := range allowedFormulaAndFunctionMetricFunctionEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}
