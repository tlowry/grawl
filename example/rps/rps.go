package main

import (
	"github.com/tlowry/grawl/browser"
	"github.com/tlowry/grawl/util"
	"log"
	"runtime/debug"
)

func main() {

	defer func() {
		if e := recover(); e != nil {
			log.Printf("Hit an error when trying to get the latest headlines from RPS: %s, %s", e, debug.Stack())
		}
	}()

	conn := browser.NewBrowser()

	page := conn.Load("rockpapershotgun.com")
	//page.Absolutify()
	page.SaveToFile("rps.html")
	posts := page.AllById(`post-*[0-9]`)

	for _, post := range posts {

		titleText := post.ByAttribute("rel", "bookmark").GetContent()

		log.Println("=====" + titleText + "=====")

		paras := post.AllByTag("p")

		for _, p := range paras {

			if !util.IsWhiteSpace(p.GetContent()) {
				log.Println("**")
				log.Println(p.GetContent())
			}
		}
		log.Println("---------------------------------------------------------------------------------------------------")
	}

}
