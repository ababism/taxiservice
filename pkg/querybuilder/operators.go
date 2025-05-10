package querybuilder

type Connector string

func AND() Connector {
	return "AND"
}
func OR() Connector {
	return "OR"
}

type Operator string

func LT() Operator {
	return "<"
}

func GT() Operator {
	return ">"
}
func GET() Operator {
	return ">="
}
func LET() Operator {
	return "<="
}

func EQ() Operator {
	return "="
}
func IS() Operator {
	return "IS"
}

func IN() Operator {
	return "IN"
}
