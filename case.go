package sqlbuilder

import (
	"strings"
)

type caseClause struct {
	clauses []sqlProvider
}

type SingleCase struct {
	when sqlProvider
	then sqlProvider
}

func When(clause sqlProvider) *SingleCase {
	sc := &SingleCase{
		when: clause,
	}

	return sc
}

func (sc *SingleCase) Then(clause sqlProvider) sqlProvider {
	sc.then = clause
	return sc
}

func (sc *SingleCase) getSQL(cache *varCache) string {
	return "WHEN " + sc.when.getSQL(cache) + " THEN " + sc.then.getSQL(cache)
}

type elseClause struct {
	clause sqlProvider
}

func Else(clause sqlProvider) sqlProvider {
	return &elseClause{clause}
}

func (e *elseClause) getSQL(cache *varCache) string {
	return "ELSE " + e.clause.getSQL(cache)
}

func Case(clauses ...sqlProvider) sqlProvider {
	return &caseClause{clauses}
}

func (c *caseClause) getSQL(cache *varCache) string {
	clauses := make([]string, len(c.clauses))

	for i, clause := range c.clauses {
		clauses[i] = clause.getSQL(cache)
	}

	return "(CASE " + strings.Join(clauses, " ") + " END)"
}
