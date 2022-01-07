package term

import (
	"fmt"
)

type GolangIntTerm int

// ReplaceVariables replaces any internal variables with the values
// to which they're bound.  Unbound variables are left as they are
func (g *GolangIntTerm) ReplaceVariables(b Bindings) Term {
	return g
}

// String provides a string representation of a term
func (g *GolangIntTerm) String() string{
	return fmt.Sprintf("%d", *g)
}

// Type indicates whether this term is an atom, number, compound, etc.
// ISO §7.2 uses the word "type" to descsribe this idea.  Constants are
// defined for each type.
func (g *GolangIntTerm) Type() int {
	return GolangInt
}

// Indicator() provides a "predicate indicator" representation of a term
func (g *GolangIntTerm) Indicator() string{
	return "golang:int"
}

// Unifies the invocant and another term in the presence of an
// environment.
// On succes, returns a new environment with additional variable
// bindings.  On failure, returns CantUnify error along with the
// original environment
func (g *GolangIntTerm) Unify(b Bindings, t Term) (Bindings, error){
	if IsVariable(t) {
		return t.Unify(b,g)
	}

	switch other := t.(type) {
	case *GolangIntTerm:
		if *other == *g {
			return b, nil
		} else {
			return nil, CantUnify
		}
	default:  //Last tier effort -- most CPU expensive
		selfString := g.String()
		otherString := t.String()
		if selfString == otherString {
			return b, nil
		} else {
			//fmt.Printf("Can not unify %q %#v with %q %#v\n", otherString, t, g, g)
			return b, CantUnify
		}
	}
}

func WrapInt(i int) Term {
	value := GolangIntTerm(i)
	return &value
}
