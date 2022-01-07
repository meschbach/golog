//Package gpi stands for Golang <-> Prolog Interface.  Provides elements to reduce the barrier of usage between the two
package gpi

import (
	"fmt"
	"github.com/meschbach/golog"
	"github.com/meschbach/golog/term"
)

//SliceAccessor is an iterator returning the set of Prolog Terms for the solution it is currently at.  Once the end has
//been reached calling any other method is considered an error
type SliceAccessor interface {
	//AsTerms provides the representation of the current position as a set of Prolog terms
	AsTerms() []term.Term
	//Next position in the slice, returning true if there are additional positions or false if the end of the slice has
	//been reached.
	Next() bool
}

//SlicePredicate adapts teh SliceAccessor into a ForeignPredicate for usage as a back trackable predicate within Golog.
type SlicePredicate struct {
	inputs []term.Term
	origin golog.Machine
	slice  SliceAccessor
}

func (i *SlicePredicate) first() golog.ForeignReturn  {
	out, err := i.Follow()
	if err == term.CantUnify {
		return golog.ForeignFail()
	} else if err != nil {
		panic(err)
	}
	return out
}

func (i *SlicePredicate) attemptPositionUnify() (term.Bindings, error) {
	terms := i.slice.AsTerms()
	if len(terms) != len(i.inputs) {
		panic(fmt.Sprintf("Expcted %d terms, got %d terms", len(i.inputs), len(terms)))
	}

	var err error
	env := i.origin.Bindings()
	for termIndex, term := range terms {
		env, err = term.Unify(env,  i.inputs[termIndex])
		if err != nil {
			return nil, err
		}
	}
	return env, nil
}

func (i *SlicePredicate) Follow() (golog.Machine, error) {
	for {
		env, err := i.attemptPositionUnify()
		if err == nil {
			more := i.slice.Next()

			next := i.origin.SetBindings(env)
			if more {
				next = next.PushDisj(i)
			}
			return next, nil
		}

		if err == term.CantUnify {
			if !i.slice.Next() {
				return nil, term.CantUnify
			}
		} else if err != nil {
			return nil, err
		}
	}
}

//IntSliceAccessor adapts a primitive int slice into Golog GolangIntTerm in the form of Index, Value
type IntSliceAccessor struct {
	elements []int
	position int
}

func (i *IntSliceAccessor) AsTerms() []term.Term {
	return []term.Term{
		term.WrapInt(i.position),
		term.WrapInt(i.elements[i.position]),
	}
}

func (i *IntSliceAccessor) Next() bool {
	i.position++
	return i.position < len(i.elements)
}

//NewIntSlicePredicate provides a golog.ForeignPredicate which provides all solutions for each call.
func NewIntSlicePredicate(args ...int) golog.ForeignPredicate {
	return func(machine golog.Machine, terms []term.Term) golog.ForeignReturn {
		it := &SlicePredicate{
			inputs: terms,
			origin: machine,
			slice: &IntSliceAccessor{
				elements: args,
				position: 0,
			},
		}

		return it.first()
	}
}
