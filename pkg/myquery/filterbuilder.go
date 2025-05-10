package myquery

import (
	"fmt"
	"strings"
)

type FilterBuilder struct {
	//baseQuery string
	query     strings.Builder
	args      []any
	lastParam int

	damaged    bool
	paramCount int

	lastArgForComparisonIsNil bool
}

func NewFilterBuilder(baseQuery string, args []any, lastParamNumber int) *FilterBuilder {
	b := strings.Builder{}
	b.WriteString(baseQuery)
	return &FilterBuilder{query: b,
		args:                      args,
		lastParam:                 lastParamNumber,
		lastArgForComparisonIsNil: true,
		damaged:                   false}
}

func (qb *FilterBuilder) Build() (query string, args []any, lastParam int) {
	return qb.query.String(), qb.args, qb.lastParam
}

func (qb *FilterBuilder) includeParams(n int) int {
	qb.lastParam += n
	qb.paramCount += n
	return qb.lastParam
}

func (qb *FilterBuilder) placeParam() int {
	qb.lastParam++
	qb.paramCount++
	return qb.lastParam
}

func (qb *FilterBuilder) Where() *FilterBuilder {
	qb.query.WriteString(" WHERE ")
	return qb
}

func (qb *FilterBuilder) ArgOpt(arg *any) *FilterBuilder {
	if arg != nil {
		qb.query.WriteString(fmt.Sprintf(" $%d", qb.placeParam()))
		qb.args = append(qb.args, *arg)
		qb.lastArgForComparisonIsNil = false
		return qb
	}

	qb.lastArgForComparisonIsNil = true
	return qb
}

func (qb *FilterBuilder) Arg(arg any) *FilterBuilder {
	qb.query.WriteString(fmt.Sprintf(" $%d", qb.placeParam()))
	qb.args = append(qb.args, arg)
	qb.lastArgForComparisonIsNil = false

	return qb
}

func (qb *FilterBuilder) LT() *FilterBuilder {
	qb.query.WriteString(" <")
	return qb
}
func (qb *FilterBuilder) GT() *FilterBuilder {
	qb.query.WriteString(" >")
	return qb
}
func (qb *FilterBuilder) GET() *FilterBuilder {
	qb.query.WriteString(" >=")
	return qb
}
func (qb *FilterBuilder) LET() *FilterBuilder {
	qb.query.WriteString(" <=")
	return qb
}

func (qb *FilterBuilder) EQ() *FilterBuilder {
	qb.query.WriteString(" =")
	return qb
}

func (qb *FilterBuilder) IS() *FilterBuilder {
	qb.query.WriteString(" IS")
	return qb
}

func (qb *FilterBuilder) IN() *FilterBuilder {
	qb.query.WriteString(" IN")
	return qb
}
func (qb *FilterBuilder) OpenBR() *FilterBuilder {
	qb.query.WriteString(" ( ")

	return qb
}
func (qb *FilterBuilder) CloseBR() *FilterBuilder {
	qb.query.WriteString(" )")
	return qb
}
func (qb *FilterBuilder) InlineWithBrackets(other FilterBuilder) *FilterBuilder {
	qb.query.WriteString(" ( ")
	qb.query.WriteString(other.query.String())
	qb.query.WriteString(" )")
	qb.args = append(qb.args, other.args)

	return qb
}

// REWORK
func (qb *FilterBuilder) WhereAllAnd(filterMap map[string]any) *FilterBuilder {
	return qb.where(filterMap, "AND")
}

func (qb *FilterBuilder) WhereAllOr(filterMap map[string]any) *FilterBuilder {
	return qb.where(filterMap, "OR")
}

func (qb *FilterBuilder) where(argsMap map[string]any, argsConnector string) *FilterBuilder {
	if len(argsMap) > 0 {
		qb.query.WriteString(" WHERE ")
	}
	connector := ""
	i := 0
	for key, val := range argsMap {
		if i == 1 {
			connector = fmt.Sprintf(" %s ", argsConnector)
		}
		qb.query.WriteString(fmt.Sprintf("%s %v = $%d", connector, key, qb.placeParam()))
		qb.args = append(qb.args, val)
		i++
	}
	return qb
}

func (qb *FilterBuilder) AddCustomQueryPart(qp string, paramAmount int) int {
	qb.query.WriteString(" " + qp)
	qb.includeParams(paramAmount)
	return qb.lastParam
}

func (qb *FilterBuilder) Limit(limit int) *FilterBuilder {

	qb.query.WriteString(fmt.Sprintf(` LIMIT $%d`, qb.placeParam()))
	qb.args = append(qb.args, limit)

	return qb
}

func (qb *FilterBuilder) Offset(offset int) *FilterBuilder {

	qb.query.WriteString(fmt.Sprintf(` OFFSET $%d`, qb.placeParam()))
	qb.args = append(qb.args, offset)

	return qb
}

func (qb *FilterBuilder) OrderBy(field string, ascending bool) *FilterBuilder {
	if ascending {
		qb.query.WriteString(fmt.Sprintf(` ORDER BY %s ASC`, field))
		return qb
	}
	qb.query.WriteString(fmt.Sprintf(` ORDER BY %s DESC`, field))
	return qb
}
