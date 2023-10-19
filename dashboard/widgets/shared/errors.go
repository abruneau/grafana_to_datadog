package shared

import "fmt"

func NoValidFormulaError(panelTitle string) error {
	return fmt.Errorf("%s no valid formula found", panelTitle)
}

func AggregationTypeError(aggregation string) error {
	return fmt.Errorf("alignement type %s not supported", aggregation)
}
