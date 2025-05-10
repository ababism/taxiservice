package querybuilder

import (
	"fmt"
	"github.com/modern-go/reflect2"
	"strings"
)

type NamedQueryBuilder struct {
	//baseQuery string
	query *strings.Builder
	args  map[string]any

	queryIncrement *strings.Builder

	incrementState bool
	//lastArgForComparisonIsNil bool

	//damaged bool
}

func (qb *NamedQueryBuilder) addQ(q string) {
	if qb.incrementState {
		if qb.queryIncrement == nil {
			qb.queryIncrement = &strings.Builder{}
		}
		qb.queryIncrement.WriteByte(' ')
		qb.queryIncrement.WriteString(q)
	} else {
		qb.query.WriteByte(' ')
		qb.query.WriteString(q)
	}
}

func (qb *NamedQueryBuilder) endIncrement() {
	if qb.incrementState {
		if qb.queryIncrement == nil {
			return
		}
		qb.query.WriteString(qb.queryIncrement.String())
		qb.queryIncrement.Reset()
		qb.incrementState = false
	}
}

func (qb *NamedQueryBuilder) dropIncrement() {
	if qb.queryIncrement != nil {
		qb.queryIncrement.Reset()
	} else {
		qb.queryIncrement = &strings.Builder{}
	}
	qb.incrementState = false
}

func (qb *NamedQueryBuilder) startIncrement() {
	if qb.incrementState {
		if qb.queryIncrement == nil {
			qb.queryIncrement = &strings.Builder{}
		} else {
			qb.queryIncrement.Reset()
		}
	}

	qb.incrementState = true
}

func NewNamed() *NamedQueryBuilder {
	var b strings.Builder
	var ib strings.Builder

	return &NamedQueryBuilder{
		query:          &b,
		queryIncrement: &ib,

		args:           make(map[string]any),
		incrementState: false,
	}
}

func NewNamedFrom(baseQuery string, args map[string]any) *NamedQueryBuilder {
	var b strings.Builder
	var ib strings.Builder
	b.WriteString(baseQuery)
	if args == nil {
		args = make(map[string]any)
	}

	return &NamedQueryBuilder{
		query:          &b,
		queryIncrement: &ib,

		args:           args,
		incrementState: false,
	}
}

func (qb *NamedQueryBuilder) Build() (query string, args map[string]any) {
	qb.dropIncrement()
	return qb.query.String(), qb.args
}

func (qb *NamedQueryBuilder) Q(q string) *NamedQueryBuilder {
	qb.addQ(q)
	return qb
}

func (qb *NamedQueryBuilder) Col(field string) *NamedQueryBuilder {
	qb.addQ(field)
	return qb
}

func (qb *NamedQueryBuilder) Table(name string) *NamedQueryBuilder {
	qb.addQ(name)
	return qb
}

func (qb *NamedQueryBuilder) CustomQ(q string, other map[string]any) *NamedQueryBuilder {
	qb.addQ(q)
	for k, v := range other {
		qb.args[k] = v
	}
	return qb
}
func (qb *NamedQueryBuilder) Returning() *NamedQueryBuilder {
	qb.addQ("RETURNING")
	return qb
}

func (qb *NamedQueryBuilder) StartOpt() *NamedQueryBuilder {
	qb.startIncrement()
	return qb
}

func (qb *NamedQueryBuilder) EndOptIf(cond func() bool) *NamedQueryBuilder {
	if cond() {
		if qb.incrementState {
			if qb.queryIncrement == nil {
				qb.queryIncrement = &strings.Builder{}
			}
			qb.query.WriteString(qb.queryIncrement.String())
			qb.queryIncrement.Reset()
			qb.incrementState = false
		}
	} else {
		qb.dropIncrement()
	}
	return qb
}

func (qb *NamedQueryBuilder) WhereOptPart() *NamedQueryBuilder {
	qb.startIncrement()
	qb.addQ("WHERE")
	return qb
}
func (qb *NamedQueryBuilder) EndWhereOpt() *NamedQueryBuilder {
	qb.dropIncrement()
	qb.TrimOpt()
	return qb
}

func (qb *NamedQueryBuilder) ArgOpt(name string, arg any) *NamedQueryBuilder {
	if !reflect2.IsNil(arg) {
		qb.addQ(fmt.Sprintf(":%s", name))
		qb.args[name] = arg

		qb.endIncrement()

		return qb
	}
	return qb
}

func (qb *NamedQueryBuilder) CompConnectorOpt(col string, op Operator, namedKey string, arg any, connector Connector) *NamedQueryBuilder {
	if !reflect2.IsNil(arg) {
		qb.addQ(fmt.Sprintf("%s %s :%s %s", col, op, namedKey, connector))
		qb.args[namedKey] = arg

		qb.endIncrement()

		return qb
	}
	return qb
}

func (qb *NamedQueryBuilder) CompOpt(col string, op Operator, namedKey string, arg any) *NamedQueryBuilder {
	if !reflect2.IsNil(arg) {
		qb.addQ(fmt.Sprintf("%s %s :%s", col, op, namedKey))
		qb.args[namedKey] = arg

		qb.endIncrement()

		return qb
	}
	return qb
}

func (qb *NamedQueryBuilder) TrimOpt() *NamedQueryBuilder {
	if qb.incrementState && qb.queryIncrement != nil {
		// Trim AND or OR in the end of queryIncrement string builder
		query := qb.queryIncrement.String()
		query = strings.TrimSuffix(query, " AND")
		query = strings.TrimSuffix(query, " OR")
		qb.queryIncrement.Reset()
		qb.queryIncrement.WriteString(query)
		return qb
	}
	// TRIM AND or OR in the end of query string builder

	query := qb.query.String()
	query = strings.TrimSuffix(query, " AND")
	query = strings.TrimSuffix(query, " OR")
	qb.query.Reset()
	qb.query.WriteString(query)
	return qb
}

//	func (qb *NamedQueryBuilder) Arg(name string, arg any) *NamedQueryBuilder {
//		qb.addQ(fmt.Sprintf(":%s", name))
//		qb.args[name] = arg
//		qb.endIncrement()
//
//		return qb
//	}
func (qb *NamedQueryBuilder) ON() *NamedQueryBuilder {
	qb.addQ("ON")
	return qb
}
func (qb *NamedQueryBuilder) LT() *NamedQueryBuilder {
	qb.addQ("<")
	return qb
}
func (qb *NamedQueryBuilder) GT() *NamedQueryBuilder {
	qb.addQ(">")
	return qb
}
func (qb *NamedQueryBuilder) GET() *NamedQueryBuilder {
	qb.addQ(">=")
	return qb
}
func (qb *NamedQueryBuilder) LET() *NamedQueryBuilder {
	qb.addQ("<=")
	return qb
}

func (qb *NamedQueryBuilder) EQ() *NamedQueryBuilder {
	qb.addQ("=")
	return qb
}

func (qb *NamedQueryBuilder) IS() *NamedQueryBuilder {
	qb.addQ("IS")
	return qb
}

func (qb *NamedQueryBuilder) IN() *NamedQueryBuilder {
	qb.addQ("IN")
	return qb
}

func (qb *NamedQueryBuilder) OpenBR() *NamedQueryBuilder {
	qb.addQ("(")
	return qb
}
func (qb *NamedQueryBuilder) CloseBR() *NamedQueryBuilder {
	qb.addQ(")")
	return qb
}
func (qb *NamedQueryBuilder) InlineWithBrackets(other *NamedQueryBuilder) *NamedQueryBuilder {
	if other == nil {
		return qb
	}
	qb.addQ("(")
	qb.addQ(other.query.String())
	qb.addQ(")")
	for k, v := range other.args {
		qb.args[k] = v
	}

	return qb
}

type ArgPair struct {
	namedAs string
	value   string
}

func (qb *NamedQueryBuilder) AllAndEQ(argsMap map[string]ArgPair) *NamedQueryBuilder {
	return qb.сonnectEQ(argsMap, "AND")
}

func (qb *NamedQueryBuilder) AllOrEQ(argsMap map[string]ArgPair) *NamedQueryBuilder {
	return qb.сonnectEQ(argsMap, "OR")
}
func (qb *NamedQueryBuilder) сonnectEQ(argsMap map[string]ArgPair, argsConnector string) *NamedQueryBuilder {
	if len(argsMap) > 0 {
		qb.addQ("WHERE")
	}

	connector := " "
	i := 0
	for key, val := range argsMap {
		if i == 1 {
			connector = fmt.Sprintf(" %s ", argsConnector)
		}

		qb.addQ(fmt.Sprintf("%s%s = :%s", connector, key, val.namedAs))
		qb.args[val.namedAs] = val.value

		i++
	}
	return qb
}

func (qb *NamedQueryBuilder) Limit(named string, limit int) *NamedQueryBuilder {
	if named == "" {
		named = "limit"
	}

	qb.addQ(fmt.Sprintf(`LIMIT :%s`, named))
	qb.args[named] = limit

	return qb
}

func (qb *NamedQueryBuilder) Offset(named *string, limit int) *NamedQueryBuilder {
	var key string
	if named == nil {
		key = "offset"
	}

	qb.addQ(fmt.Sprintf(`OFFSET :%s`, key))
	qb.args[key] = limit

	return qb
}

func (qb *NamedQueryBuilder) OrderBy(col string, ascending bool) *NamedQueryBuilder {
	if ascending {
		qb.addQ(fmt.Sprintf(`ORDER BY %s ASC`, col))
		return qb
	}
	qb.addQ(fmt.Sprintf(`ORDER BY %s DESC`, col))
	return qb
}
