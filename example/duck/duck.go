package main

import (
	"github.com/tlowry/grawl/browser"
	"github.com/tlowry/grawl/element"
	"log"
	"runtime/debug"
	"strings"
)

/*
	TODO :
	handle text differently (not nodes)
	Add more examples
	Add cookie manipulation in browser
	Add godoc
	Add test/bench (maybe problem on windows)
	Add a relative url conversion option for file output
	look at TODO comments
	Add javascript support
*/

func main() {

	defer func() {
		if e := recover(); e != nil {
			log.Printf("Hit an error when trying to get results : %s, %s", e, debug.Stack())
		}
	}()

	conn := browser.NewBrowser()
	conn.SetUserAgent("Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.137 Safari/537.36")

	page := conn.Load("https://duckduckgo.com")

	page.SaveToFile("b4.html")

	formElem := page.ById("search_form_homepage")

	failOnNil("Couldn't find the search form", formElem)

	form := formElem.(*element.Form)

	form.SetField("q", "hello")

	page = conn.SubmitForm(form)
	page.SaveToFile("afterForm.html")

	results := page.AllByAttribute("class", "results_links results_links_deep web-result")

	var elem element.Element = nil
	var snipText string

	log.Println("found ", len(results), " results")

	// Look through the results page for anything to do with "Kitty"
	for _, el := range results {
		snippet := el.ByAttribute("class", "snippet")
		snipText = snippet.GetContent()
		if strings.Contains(snipText, "Kitty") {
			elem = el
			break
		}
	}

	failOnNil("Couldn't find anything related to Kitty", elem)

	link := elem.ByTag("a")
	url := link.GetAttribute("href")
	log.Println("Found a result at " + url + "\n" + "\"" + snipText + "\"")
	resultPage := conn.Load(url)

	// Save the page to file
	resultPage.SaveToFile("out.html")
}

// Quit the app with the provided error string message if iface is null
func failOnNil(msg string, iface interface{}) {
	if iface == nil {
		log.Fatal(msg)
	}
}
