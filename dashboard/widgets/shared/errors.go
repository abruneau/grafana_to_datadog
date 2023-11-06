package shared

import "fmt"

func NoValidFormulaError(panelTitle string) error {
	return fmt.Errorf("%s no valid formula found", panelTitle)
}

func AggregationTypeError(aggregation string, a ...any) error {
	return fmt.Errorf("alignement type %s not supported %s", aggregation, a)
}

func TransformationTypeError(aggregation any, a ...any) error {
	return fmt.Errorf("transformation type %s not supported %s", aggregation, a)
}
