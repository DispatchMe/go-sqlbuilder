package sqlbuilder

// Add a HAVING clause to your query with one or more constraints (either Expr instances or And/Or functions)
func (q *Query) Having(constraints ...sqlProvider) *Query {
	if q.having == nil {
		q.having = new(constraint)
		q.having.gate = gate_and
	}

	q.having.children = append(q.having.children, constraints...)

	return q
}
