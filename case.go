package sqlbuilder

import (
	"strings"
)

type caseClause struct {
	clauses []SQLProvider
}

type SingleCase struct {
	when SQLProvider
	then SQLProvider
}

func When(clause SQLProvider) *SingleCase {
	sc := &SingleCase{
		when: clause,
	}

	return sc
}

func (sc *SingleCase) Then(clause SQLProvider) SQLProvider {
	sc.then = clause
	return sc
}

func (sc *SingleCase) getSQL(cache *varCache) string {
	return "WHEN " + sc.when.getSQL(cache) + " THEN " + sc.then.getSQL(cache)
}

type elseClause struct {
	clause SQLProvider
}

func Else(clause SQLProvider) SQLProvider {
	return &elseClause{clause}
}

func (e *elseClause) getSQL(cache *varCache) string {
	return "ELSE " + e.clause.getSQL(cache)
}

func Case(clauses ...SQLProvider) SQLProvider {
	return &caseClause{clauses}
}

func (c *caseClause) getSQL(cache *varCache) string {
	clauses := make([]string, len(c.clauses))

	for i, clause := range c.clauses {
		clauses[i] = clause.getSQL(cache)
	}

	return "(CASE " + strings.Join(clauses, " ") + " END)"
}
