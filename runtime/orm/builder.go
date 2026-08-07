package orm

func (m *Model) Table(name string) *Query {

	return &Query{
		model: m,
		table: name,
	}
}

func (q *Query) Where(
	query string,
	args ...any,
) *Query {

	q.where = query

	q.args = args

	return q
}