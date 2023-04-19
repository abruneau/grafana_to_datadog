package cloudwatch

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

func NewCloudwatchFormula(expression string) (*datadogV1.WidgetFormula, error) {
	formula := cloudwatchFormula{expression: expression}
	return formula.build()
}

type cloudwatchFormula struct {
	expression string
	parts      []string
}

func (f *cloudwatchFormula) build() (*datadogV1.WidgetFormula, error) {
	err := f.splitExpression()
	if err != nil {
		return nil, err
	}

	query := strings.Join(f.parts, "")
	return datadogV1.NewWidgetFormula(query), nil
}

func (f *cloudwatchFormula) splitExpression() error {
	f.parts = []string{}
	exp := regexp.MustCompile(`[A-Z_]+\(.*?\)`)
	matches := exp.FindAllStringIndex(f.expression, -1)

	// check if there is no match
	if len(matches) == 0 {
		f.parts = append(f.parts, f.expression)
		return nil
	}

	// check if the first match is not the begining of the expression
	if matches[0][0] > 0 {
		f.parts = append(f.parts, f.expression[:matches[0][0]])
	}

	for i, m := range matches {
		exp, err := f.parseFunction(f.expression[m[0]:m[1]])
		if err != nil {
			return err
		}
		f.parts = append(f.parts, exp)

		// check if there are elements between this match and the next one
		if i < len(matches)-1 && matches[i+1][0]-1 > m[1] {
			f.parts = append(f.parts, f.expression[m[1]:matches[i+1][0]])
		}
	}

	if len(f.expression) > matches[len(matches)-1][1] {
		f.parts = append(f.parts, f.expression[matches[len(matches)-1][1]:])
	}
	return nil
}

func (f *cloudwatchFormula) parseFunction(part string) (string, error) {
	exp := regexp.MustCompile(`[A-Z_]{2,}`)
	operator := exp.FindStringSubmatch(part)[0]

	switch operator {
	case "SUM":
		return f.parseSumFunction(part), nil
	default:
		return "", fmt.Errorf("unknown operator %s", operator)
	}
}

func (f *cloudwatchFormula) parseSumFunction(part string) string {
	return fmt.Sprintf("(%s)", strings.ReplaceAll(part[5:len(part)-2], ",", "+"))
}
