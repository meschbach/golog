//Package gpi stands for Golang <-> Prolog Interface.  Provides elements to reduce the barrier of usage between the two
package gpi

import (
	"fmt"
	"github.com/meschbach/golog"
	"github.com/meschbach/golog/term"
)

type IntSlicePredicate struct {
	inputs   []term.Term
	origin   golog.Machine
	elements []int
	index    int
}

func (i *IntSlicePredicate) Follow() (golog.Machine, error) {
	currentIndex := i.index
	i.index++

	value := i.elements[currentIndex]

	indexTerm := term.NewInt64(int64(currentIndex))
	valueTerm := term.NewInt64(int64(value))

	terms := []term.Term{
		indexTerm,
		valueTerm,
	}
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
	if i.index < len(i.elements) {
		next = next.PushDisj(i)
	}
	return next, nil

}

func NewIntSlicePredicate(args ...int) golog.ForeignPredicate {
	return func(machine golog.Machine, terms []term.Term) golog.ForeignReturn {
		it := &IntSlicePredicate{
			inputs:   terms,
			origin:   machine,
			elements: args,
			index:    0,
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
