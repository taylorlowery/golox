package golox

type Binary struct {
	Expr left
	Token operator
	Expr right
}

type Grouping struct {
	Expr expression
}

type Literal struct {
	Object value
}

type Unary struct {
	Token operator
	Expr right
}

