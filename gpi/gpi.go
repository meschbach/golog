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

func (i *SlicePredicate) Follow() (golog.Machine, error) {
	terms := i.slice.AsTerms()
	more := i.slice.Next()
	if len(terms) != len(i.inputs) {
		panic(fmt.Sprintf("Expcted %d terms, got %d terms", len(i.inputs), len(terms)))
	}

	var err error
	env := i.origin.Bindings()
	for termIndex := 0; termIndex < len(terms); termIndex++ {
		env, err = i.inputs[termIndex].Unify(env, terms[termIndex])
		if err != nil {
			return nil, err
		}
	}

	next := i.origin.SetBindings(env)
	if more {
		next = next.PushDisj(i)
	}
	return next, nil
}

type IntSliceAccessor struct {
	elements []int
	position int
}

func (i *IntSliceAccessor) AsTerms() []term.Term {
	return []term.Term{
		term.NewInt64(int64(i.position)),
		term.NewInt64(int64(i.elements[i.position])),
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

		m, err := it.Follow()
		//Filter through each possibility until we can resolve the constraints
		for ; err == term.CantUnify; m, err = it.Follow() {

		}
		if err != nil {
			//TODO: Should not eat error
			//panic(err)
			return golog.ForeignFail()
		}
		return m
	}
}
