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
	for {
		if err := update(); err != nil {
			ui.ErrorMessage = fmt.Sprintf("Error while fetching: %s", err)
			ui.FooterMessage = fmt.Sprintf("Last attempted update at %s", time.Now())
		} else {
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
	feed, err := fetchFeed(fetchURL)

	if err != nil {
		return err
	}
	for i := len(feed.Entries) - 1; i >= 0; i-- {
		e := feed.Entries[i]
		entries.Add(e)
	}
	entries.Sort()
	return nil
}
