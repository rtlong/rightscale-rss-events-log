package ui

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/dustin/go-humanize"
	textwrap "github.com/kr/text"
	"github.com/nsf/termbox-go"
)

type UI struct {
	Collection    EntrySet
	HeaderMessage string
	FooterMessage string
	ErrorMessage  string

	DefaultBG termbox.Attribute
	DefaultFG termbox.Attribute

	Events chan Event

	drawVent chan bool
	wg       sync.WaitGroup

	width  int
	height int

	initted bool
}

func (ui *UI) Init() {
	if !ui.initted {
		ui.drawVent = make(chan bool, 0)
		ui.Events = make(chan Event, 0)
		ui.Collection = NewEntrySet()
		ui.initted = true
	}
}

func (ui *UI) Start() (err error) {
	ui.Init()
	err = termbox.Init()
	if err != nil {
		return
	}
	termbox.SetInputMode(termbox.InputEsc)
	ui.width, ui.height = termbox.Size()

	// ui.wg.Add(1)
	go ui.drawLoop()

	go ui.eventLoop()

	// Initial draw
	ui.Redraw()

	return
}

func (ui *UI) Stop() {
	// close(ui.drawVent)
	ui.wg.Wait()

	termbox.Close()

	close(ui.Events)
}

func (ui *UI) Redraw() {
	// log.Println("Redraw:start")
	// log.Printf("Redraw:start ui.drawVent=%#v\n", ui.drawVent)
	ui.drawVent <- true
	// log.Println("Redraw:end")
}

func (ui *UI) drawLoop() {
	var y int
	// log.Printf("drawLoop:start ui.drawVent=%#v\n", ui.drawVent)
	for _ = range ui.drawVent {
		// log.Println("drawLoop:loop")
		termbox.Clear(ui.DefaultFG, ui.DefaultBG)
		termbox.HideCursor()

		y = ui.printHeader(0, 0)
		y = ui.printCollection(0, y, ui.height-y-5)
		y += 1
		y = ui.printFooter(0, y)
		y = ui.printError(0, y)

		termbox.Flush()
		// log.Println("drawLoop:loop:end")
	}
}

func (ui *UI) eventLoop() {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Ch == 0 {
				switch ev.Key {
				case termbox.KeyEsc, termbox.KeyCtrlC:
					ui.Events <- Event{Type: EventQuit}
				}
			} else {
				switch ev.Ch {
				case 'q':
					ui.Events <- Event{Type: EventQuit}
				}
			}

		case termbox.EventResize:
			ui.width = ev.Width
			ui.height = ev.Height
			ui.Redraw()

		case termbox.EventError:
			ui.Events <- Event{Type: EventError, Err: ev.Err}
		}
	}
}

func (ui *UI) printHeader(x int, y int) int {
	if ui.HeaderMessage != "" {
		_, y = ui.printText(x, y, ui.HeaderMessage, termbox.ColorCyan, ui.DefaultBG)
		y += 1
	}
	return y
}

func (ui *UI) printFooter(x, y int) int {
	if ui.FooterMessage != "" {
		_, y = ui.printText(x, y, ui.FooterMessage, termbox.ColorCyan, ui.DefaultBG)
		y += 1
	}
	return y
}

func (ui *UI) printError(x, y int) int {
	if ui.ErrorMessage != "" {
		lines := ui.wrapString(x, ui.ErrorMessage)
		for _, line := range lines {
			_, y = ui.printText(x, y, line, termbox.ColorRed, ui.DefaultBG)
			y += 1
		}
	}
	return y
}

func (ui *UI) printCollection(originX, originY, maxY int) (y int) {
	x, y := originX, maxY
	// totalRows := maxY - originY
	var entries []Entry = ui.Collection.Entries

	// if l := ui.Collection.Len(); l > totalRows {
	// entries = ui.Collection.Entries[l-totalRows : l]
	// }
	//
	secondColumnX := originX + 25

	for i := len(entries) - 1; y >= originY && i >= 0; i-- {
		item := entries[i]
		// log.Printf("printCollection:entry item=%#v\n", item)
		time := fmt.Sprintf("[%s]", humanize.Time(item.Updated()))

		titleLines := ui.wrapString(secondColumnX, item.Title())
		y -= len(titleLines)

		x, y = ui.printText(x, y, time, ui.DefaultFG, ui.DefaultBG)

		_y := y
		for _, line := range titleLines {
			x, _y = ui.printText(secondColumnX, _y, line, ui.DefaultFG, ui.DefaultBG)
			_y += 1
		}

		x = originX
	}

	return maxY
}

var (
	sp = []byte{' '}
)

// quantity of rows needed for wrapped string given starting X position
func (ui *UI) wrapString(originX int, text string) (lines []string) {
	w := ui.wrapWidth(originX)

	if len(text) <= w {
		lines = []string{text}
		return
	}

	words := bytes.Split(bytes.Replace(bytes.TrimSpace([]byte(text)), []byte{'\n'}, sp, -1), sp)
	wrappedLines := textwrap.WrapWords(words, 1, w, 1e5)
	lines = make([]string, len(wrappedLines))
	for i, words := range wrappedLines {
		lines[i] = string(bytes.Join(words, sp))
	}
	return
}

func (ui *UI) wrapWidth(originX int) int {
	return ui.width - originX - 1
}

func (ui *UI) printText(originX, originY int, text string, fg, bg termbox.Attribute) (x, y int) {
	x, y = originX, originY

	for _, r := range text {
		termbox.SetCell(x, y, r, fg, bg)
		x++
	}

	return
}
