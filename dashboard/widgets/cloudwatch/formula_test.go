package cloudwatch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitExpression(t *testing.T) {
	formula := new(cloudwatchFormula)

	var tests = []struct {
		input    string
		expected []string
	}{
		{
			"SUM([e400,e500])*100/countt", []string{"(e400+e500)", "*100/countt"},
		},
		{
			"5*SUM([e400,e500])*100/SUM([e400,e500])*countt", []string{"5*", "(e400+e500)", "*100/", "(e400+e500)", "*countt"},
		},
	}

	for _, test := range tests {
		formula.expression = test.input
		formula.splitExpression()
		assert.Equal(t, len(formula.parts), len(test.expected))
		for i, p := range formula.parts {
			assert.Equal(t, test.expected[i], p)
		}
	}
}
