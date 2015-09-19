package sqlbuilder

// Add a WHERE clause
func (q *Query) Where(constraints ...sqlProvider) *Query {
	if q.where == nil {
		q.where = new(constraint)
		q.where.gate = GATE_AND
	}

	q.where.children = append(q.where.children, constraints...)

	return q
}
