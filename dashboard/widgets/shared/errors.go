package shared

import "fmt"

func NoValidFormulaError(panelTitle string) error {
	return fmt.Errorf("%s no valid formula found", panelTitle)
}
