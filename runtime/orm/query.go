package orm

type Query struct {
	model *Model

	table string

	where string

	args []any
}