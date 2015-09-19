package sqlbuilder

import (
	"strings"
)

type constraint struct {
	gate            int
	children        []sqlProvider
	field           string
	comparator      string
	value           interface{}
	bypassStatement bool
}

func (c *constraint) getSQL(cache *varCache) string {
	if len(c.children) > 0 {
		compiled := make([]string, len(c.children))
		for i, cstr := range c.children {
			compiled[i] = cstr.getSQL(cache)
		}

		var gate string
		if c.gate == GATE_AND {
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
		var key string
		if c.bypassStatement {
			var ok bool
			key, ok = c.value.(string)
			if !ok {
				panic("Cannot generate SQL for constraint outside of a statement when the constraint value is not a string!")
			}
		} else {
			key = cache.add(c.value)
		}

		return c.field + " " + c.comparator + " " + key
	}
}

func (c *constraint) addChild(child sqlProvider) {
	c.children = append(c.children, child)
}

// This is used externally to construct query constraints
type Expr struct {
	Field      string
	Comparator string
	Value      interface{}
}

func (expr Expr) getSQL(cache *varCache) string {
	key := cache.add(expr.Value)
	return expr.Field + " " + expr.Comparator + " " + key
}

func And(constraints ...sqlProvider) *constraint {
	return &constraint{
		gate:     GATE_AND,
		children: constraints,
	}
}

func Or(constraints ...sqlProvider) *constraint {
	return &constraint{
		gate:     GATE_OR,
		children: constraints,
	}
}
