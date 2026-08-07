package orm

import "fmt"

func (q *Query) First() {

	sql := fmt.Sprintf(
		"SELECT * FROM %s WHERE %s LIMIT 1",
		q.table,
		q.where,
	)

	fmt.Println(sql)
}