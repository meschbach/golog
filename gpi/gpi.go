//Package gpi stands for Golang <-> Prolog Interface.  Provides elements to reduce the barrier of usage between the two
package gpi

import (
	"fmt"
	"github.com/meschbach/golog"
	"github.com/meschbach/golog/term"
)

type SliceAccessor interface {
	AsTerms() []term.Term
	Next() bool
}

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
