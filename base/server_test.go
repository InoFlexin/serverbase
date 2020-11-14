package base

import "testing"

func TestHelloPuppy(t *testing.T) {
	expected := "Hello, Cute puppy :)"
	if actual := Dummy(); actual != expected {
		t.Errorf("Expect - %v, but got - %v", expected, actual)
	}
}
