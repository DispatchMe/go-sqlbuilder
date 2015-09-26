package sqlbuilder

import (
	"strings"
)

type constraint struct {
	gate     int
	children []sqlProvider
	expr     sqlProvider
}

func (c *constraint) getSQL(cache *varCache) string {
	if len(c.children) > 0 {
		compiled := make([]string, len(c.children))
		for i, cstr := range c.children {
			compiled[i] = cstr.getSQL(cache)
		}

		var gate string
		if c.gate == gate_and {
			gate = " AND "
		} else {
			gate = " OR "
		}

		prefix := "("
		suffix := ")"
		if len(c.children) == 1 {
			prefix = ""
			suffix = ""
		}

		return prefix + strings.Join(compiled, gate) + suffix
	} else {
		return c.expr.getSQL(cache)
	}
}

func (c *constraint) addChild(child sqlProvider) {
	c.children = append(c.children, child)
}

// Used within a WHERE or HAVING clause to group Expr instances, or nested And/Or functions, in an AND logical gate (all constraints within this function must be true)
func And(constraints ...sqlProvider) *constraint {
	return &constraint{
		gate:     gate_and,
		children: constraints,
	}
}

// Used within a WHERE or HAVING clause to group Expr instances, or nested And/Or functions, in an OR logical gate (at least one of the constraints within this function must be true)
func Or(constraints ...sqlProvider) *constraint {
	return &constraint{
		gate:     gate_or,
		children: constraints,
	}
}
