package term

import "testing"

func TestGolangIntTerm_Unify_WithGolangIntTerm(t *testing.T) {
	t1 := WrapInt(1)
	t2 := WrapInt(1)

	bindings := NewBindings()
	_, err := t1.Unify(bindings,t2)
	if err == CantUnify {
		t.Errorf("Expected unification of %s to %s....did not", t1, t2)
	}
}