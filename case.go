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

func (sc *SingleCase) GetSQL(cache *VarCache) string {
	return "WHEN " + sc.when.GetSQL(cache) + " THEN " + sc.then.GetSQL(cache)
}

type elseClause struct {
	clause SQLProvider
}

func Else(clause SQLProvider) SQLProvider {
	return &elseClause{clause}
}

func (e *elseClause) GetSQL(cache *VarCache) string {
	return "ELSE " + e.clause.GetSQL(cache)
}

func Case(clauses ...SQLProvider) SQLProvider {
	return &caseClause{clauses}
}

func (c *caseClause) GetSQL(cache *VarCache) string {
	clauses := make([]string, len(c.clauses))

	for i, clause := range c.clauses {
		clauses[i] = clause.GetSQL(cache)
	}

	return "(CASE " + strings.Join(clauses, " ") + " END)"
}
