package ev

import "testing"

type TestToEvFoo struct {
	I int
	S string
}

func TestToEv(t *testing.T) {
	ev := ToEv(TestToEvFoo{
		I: 42,
		S: "foo",
	})
	if ev.Name != "TestToEvFoo" {
		t.Fatal()
	}
	if len(ev.Attrs) != 2 {
		t.Fatal()
	}
	if ev.Attrs[0].Name != "I" {
		t.Fatal()
	}
	if ev.Attrs[0].Value != 42 {
		t.Fatal()
	}
}
