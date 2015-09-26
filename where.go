package sqlbuilder

// Add a WHERE clause to your query with one or more constraints (either Expr instances or And/Or functions)
func (q *Query) Where(constraints ...sqlProvider) *Query {
	if q.where == nil {
		q.where = new(constraint)
		q.where.gate = gate_and
	}

	q.where.children = append(q.where.children, constraints...)

	return q
}
