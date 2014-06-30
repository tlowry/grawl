package main

import (
	"github.com/tlowry/grawl/browser"
	"log"
)

func main() {

	defer func() {
		if e := recover(); e != nil {
			log.Println("Hit an error when trying to get the latest headlines from gnews: ", e)
		}
	}()

	conn := browser.NewBrowser()
	page := conn.Load("https://news.google.com/")

	topSection := page.ByClass("section-stream-content*")

	stories := topSection.AllByClass("titletext*")

	for _, story := range stories {
		log.Println("Got story " + story.GetContent())
	}

	page.SaveToFile("out.html")

}
