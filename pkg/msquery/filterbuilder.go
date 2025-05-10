package msquery

import (
	"fmt"
	"strings"
)

type FilterBuilder struct {
	baseQuery string
	query     strings.Builder
	args      []any
	lastParam int
}

func (qb *FilterBuilder) Build() (query string, args []any, lastParam int) {
	return qb.baseQuery + qb.query.String(), qb.args, qb.lastParam
}

func NewFilterBuilder(baseQuery string, lastParamNumber int) *FilterBuilder {
	return &FilterBuilder{baseQuery: baseQuery,
		args:      make([]any, 0),
		lastParam: lastParamNumber}
}

func (qb *FilterBuilder) WhereAnd(filterMap map[string]any) *FilterBuilder {
	return qb.where(filterMap, "AND")
}

func (qb *FilterBuilder) WhereOr(filterMap map[string]any) *FilterBuilder {
	return qb.where(filterMap, "OR")
}

func (qb *FilterBuilder) placeParam() int {
	qb.lastParam++
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
