package ui

import (
	"sort"
	"time"
)

type Entry interface {
	ID() string
	Title() string
	Updated() time.Time
}

type EntrySet struct {
	Entries []Entry
}

func NewEntrySet() (s EntrySet) {
	s.Entries = make([]Entry, 0)
	return
}

func (s *EntrySet) Add(newEntry Entry) {
	// log.Printf("EntrySet.Add(%#v)\n", newEntry)
	for i, e := range s.Entries {
		if e.ID() == newEntry.ID() {
			s.Entries[i] = newEntry
			return
		}
	}
	s.Entries = append(s.Entries, newEntry)
}

func (s *EntrySet) Sort() {
	sort.Stable(s)
}

func (s *EntrySet) Len() int {
	return len(s.Entries)
}

func (s *EntrySet) Less(i, j int) bool {
	ei, ej := s.Entries[i], s.Entries[j]
	return ei.Updated().Before(ej.Updated())
}

func (s *EntrySet) Swap(i, j int) {
	s.Entries[i], s.Entries[j] = s.Entries[j], s.Entries[i]
}
