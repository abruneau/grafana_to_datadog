package prometheus

import (
	"fmt"
	"strings"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

func (s *Structure) transvers(refId string, level int) (f string, q []struct {
	Name        string
	Query       string
	Aggregation datadogV1.FormulaAndFunctionMetricAggregation
}, err error) {
	if s.Agg != "" && len(s.Args) == 2 {
		f1, q1, _ := s.Args[0].transvers(refId, level+10)
		f2, q2, _ := s.Args[1].transvers(refId, level+20)
		f = fmt.Sprintf("%s %s %s", f1, s.Agg, f2)
		q = append(q, q1...)
		q = append(q, q2...)
		return
	}

	if s.Number != "" {
		f = s.Number
	}

	if s.Function != "" && len(s.Args) > 0 {
		var formulas []string
		for i, a := range s.Args {
			fchild, qchild, _ := a.transvers(refId, level+10*(i+1))
			formulas = append(formulas, fchild)
			q = append(q, qchild...)
		}
		f = fmt.Sprintf("%s(%s)", s.Function, strings.Join(formulas, ", "))
	}

	if s.Parsed != "" {
		var id = refId
		if level > 0 {
			id = fmt.Sprintf("%s%v", refId, level)
		}
		if s.Function != "" {
			f = fmt.Sprintf("%s(%s)", s.Function, id)
		} else {
			f = id
		}
		q = append(q, struct {
			Name        string
			Query       string
			Aggregation datadogV1.FormulaAndFunctionMetricAggregation
		}{id, s.Parsed, s.Agg})
	}
	return
}
