package gpi

import (
	"github.com/meschbach/golog"
	"github.com/meschbach/golog/term"
	"testing"
)

func assertInt(t *testing.T, binding term.Bindings, name string, value int) {
	termBinding, err := binding.ByName(name)
	if err != nil {
		t.Errorf("Binding problem for %q: %s", name, err)
		return
	}
	var ok bool
	var numberTerm term.Number
	if numberTerm, ok = termBinding.(term.Number); !ok {
		t.Errorf("%q is not a number", name)
		return
	}
	bigInt, fits := numberTerm.LosslessInt()
	if !fits {
		t.Errorf("%q is not an integer, got %s", name, numberTerm)
		return
	}
	actualValue := bigInt.Int64()
	if actualValue != int64(value) {
		t.Errorf("Expected %q to be %d, got %d", name, value, actualValue)
	}
}

func TestSimpleBacktracking(t *testing.T) {
	example := []int{0, 1, 2, 3, 4, 10}
	m := golog.NewMachine().RegisterForeign(map[string]golog.ForeignPredicate{
		"int_slice/2": NewIntSlicePredicate(example...),
	})
	resultSet := m.ProveAll("int_slice(I,N).")
	for i, binding := range resultSet {
		assertInt(t, binding, "I", i)
		assertInt(t, binding, "N", example[i])
	}
	if len(resultSet) < len(example) {
		t.Errorf("Expected %d results, received %d", len(example), len(resultSet))
	}
}

func TestMatchingElement(t *testing.T) {
	example := []int{0, 1, 2, 3, 4, 10}
	m := golog.NewMachine().RegisterForeign(map[string]golog.ForeignPredicate{
		"int_slice/2": NewIntSlicePredicate(example...),
	})
	resultSet := m.ProveAll("int_slice(I,4).")
	for _, binding := range resultSet {
		assertInt(t, binding, "I", 4)
	}
	if len(resultSet) != 1 {
		t.Errorf("Expected %d results, received %d", 1, len(resultSet))
	}
}

func TestMatchingIndex(t *testing.T) {
	example := []int{0, 1, 2, 3, 4, 10}
	m := golog.NewMachine().RegisterForeign(map[string]golog.ForeignPredicate{
		"int_slice/2": NewIntSlicePredicate(example...),
	})
	resultSet := m.ProveAll("int_slice(3,I).")
	for _, binding := range resultSet {
		assertInt(t, binding, "I", 3)
	}
	if len(resultSet) != 1 {
		t.Errorf("Expected %d results, received %d", 1, len(resultSet))
	}
}
