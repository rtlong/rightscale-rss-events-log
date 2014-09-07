package main

import "time"

type Feed struct {
	Title   string    `xml:"title"`
	Updated time.Time `xml:"updated"`
	Entries []Entry   `xml:"entry"`
}

type Entry struct {
	XML_ID      string    `xml:"id"`
	XML_Title   string    `xml:"title"`
	XML_Updated time.Time `xml:"updated"`
}

func (e Entry) Title() string {
	return e.XML_Title
}

func (e Entry) Updated() time.Time {
	return e.XML_Updated
}

func (e Entry) ID() string {
	return e.XML_ID + e.XML_Title
}
