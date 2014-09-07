package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log_ui "github.com/rtlong/rightscale-rss-events-log/ui"
)

var (
	entries    log_ui.EntrySet
	fetchURL   string
	ui         log_ui.UI
	httpClient *http.Client
)

func main() {
	flag.Parse()
	if fetchURL = flag.Arg(0); fetchURL == "" {
		log.Fatal("Supply the RightScale events RSS feed URL as the only argument")
	}

	// logFile, err := os.Create("/tmp/rs_tail_log")
	// if err != nil {
	// log.Fatal("error opening log file:", err)
	// }
	// log.SetOutput(logFile)
	// defer logFile.Close()
	// fmt.Fprintln(logFile, "\n\n")
	// log.Println("main:start")

	entries = log_ui.NewEntrySet()

	go handleSigs()

	ui := log_ui.UI{
		FooterMessage: "Loading...",
	}

	if err := ui.Start(); err != nil {
		panic(err)
	}

	go updateLoop(&ui)

	for e := range ui.Events {
		// log.Printf("ui.Events:loop: %#v\n", e)
		switch e.Type {
		case log_ui.EventError:
			panic(e.Err)
		case log_ui.EventQuit:
			ui.Stop()
			break
		}
	}
}

func handleSigs() {
	sigchan := make(chan os.Signal, 3)
	signal.Notify(sigchan, os.Interrupt, os.Kill, syscall.SIGTERM)

	<-sigchan
	ui.Stop()
	os.Exit(0)
}

func updateLoop(ui *log_ui.UI) {
	// log.Printf("updateLoop ui=%#v\n", ui)
	for {
		// log.Println("updateLoop:loop")
		if err := update(); err != nil {
			// log.Printf("updateLoop:loop:error %#v\n", err)
			ui.ErrorMessage = fmt.Sprintf("Error while fetching: %s", err)
			ui.FooterMessage = fmt.Sprintf("Last attempted update at %s", time.Now())
		} else {
			// log.Printf("updateLoop:loop: %#v\n", entries)
			ui.Collection = entries
			ui.ErrorMessage = ""
			ui.FooterMessage = fmt.Sprintf("Last updated at %s", time.Now())
		}
		ui.Redraw()
		time.Sleep(time.Duration(3) * time.Second)
	}
}

func fetchFeed(url string) (feed Feed, err error) {
	httpClient = &http.Client{
		Timeout: time.Duration(5) * time.Second,
	}
	resp, err := httpClient.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	feed = Feed{}

	b := bytes.Buffer{}
	b.ReadFrom(resp.Body)

	err = xml.Unmarshal(b.Bytes(), &feed)
	if err != nil {
		return
	}
	return
}

func update() error {
	// log.Println("update")
	feed, err := fetchFeed(fetchURL)
	// log.Printf("update:after feed=%#v err=%#v\n", feed, err)

	if err != nil {
		return err
	}
	for i := len(feed.Entries) - 1; i >= 0; i-- {
		e := feed.Entries[i]
		// log.Printf("Entry: id=%s title=%#v\n", e.ID(), e.Title())
		entries.Add(e)
	}
	entries.Sort()
	return nil
}
