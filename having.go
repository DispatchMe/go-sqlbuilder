package sqlbuilder

// Add a HAVING clause to your query
func (q *Query) Having(constraints ...sqlProvider) *Query {
	if q.having == nil {
		q.having = new(constraint)
		q.having.gate = GATE_AND
	}

	q.having.children = append(q.having.children, constraints...)

	return q
}
