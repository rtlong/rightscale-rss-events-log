package ui

import (
	"testing"
	"time"
)

type Foo struct {
	id string
}

func (f Foo) ID() string {
	return f.id
}

func (f Foo) Title() string {
	return ""
}

func (f Foo) Updated() time.Time {
	return time.Now()
}

func TestAddIdentityChecking(t *testing.T) {
	s := NewEntrySet()
	item1 := &Foo{id: "id"}
	s.Add(item1)
	item2 := &Foo{id: "id"}
	s.Add(item2)

	e := 1
	if a := len(s.Entries); a != e {
		t.Errorf("Expected len(EntrySet.Entries) == %d, got %d\n", e, a)
	}
}
