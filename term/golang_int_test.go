package term

import (
	"fmt"
	"testing"
)

func TestGolangIntTerm_Unify_WithGolangIntTerm(t *testing.T) {
	t.Parallel()
	t1 := WrapInt(1)
	t2 := WrapInt(1)

	bindings := NewBindings()
	_, err := t1.Unify(bindings,t2)
	if err == CantUnify {
		t.Errorf("Expected unification of %s to %s....did not", t1, t2)
	}
}

func TestGolangIntTerm_Unify_Atom(t *testing.T) {
	t.Parallel()
	t1 := WrapInt(314)
	t2 := NewAtom("314")

	fmt.Printf("T2 string: %q\n",t2)
	bindings := NewBindings()
	_, err := t1.Unify(bindings,t2)
	if err != CantUnify {
		t.Errorf("Expected unification to fail of %s to %s....did not", t1, t2)
	}
}