package main

import (
	"github.com/tlowryie/grawl/browser"
	"github.com/tlowryie/grawl/element"
	"log"
)

func main() {

	defer func() {
		if e := recover(); e != nil {
			log.Println("Hit an error when trying to get the latest headlines from RPS: ", e)
		}
	}()

	conn := browser.NewBrowser()

	page := conn.Load("http://www.google.ie/")
	page.SaveToFile("b4.html")

	form := page.ByAttribute("name", "f").(*element.Form)

	form.SetField("q", "Hello")

	page = conn.SubmitForm(form)

	page.SaveToFile("out.html")

}
