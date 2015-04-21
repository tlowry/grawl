package element

import (
	"bufio"
	"io"
	"net/http"
	"os"
)

type Page struct {
	Document *http.Response
	root     Element
	url      string
}

func NewPage() *Page {
	p := Page{}
	p.Document = nil
	return &p
}

// Build a Page from a http response
func ParseResp(resp *http.Response) *Page {

	r := bufio.NewReader(resp.Body)
	defer resp.Body.Close()

	p := ParseBody(r)

	p.Document = resp
	return p

}

func ParseBody(r io.Reader) *Page {
	parser := NewParser()
	return parser.ParsePage(r)
}

func (p *Page) GetUrl() string {
	return p.url
}

func (p *Page) SetUrl(url string) {
	p.url = url
}

/*
	Convert all relative links on this page to absolute links
	(useful when saving a file to disk for later viewing)
*/
func (p *Page) Absolutify() {
	if p.root != nil {
		allLinks := p.root.AllByTag("a")

		for _, link := range allLinks {
			url := link.GetAttribute("href")
			if url[0] == '/' {
				link.SetAttribute("href", p.GetUrl()+"/")
			}
		}
	}
}

/*
	Blocking call to save this pages html content to a file
*/
func (p *Page) SaveToFile(fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	} else {
		defer f.Close()
		err := ElementToFile(p.root, f)
		if err != nil {
			panic(err)
		}
	}

}

/*
	Saves a textual markup representation of an
	Element and all of it's child elements to a file
	TODO this is recursive (careful of stack overflows)
*/
func ElementToFile(e Element, out *os.File) (err error) {
	var tagContent string

	tagContent = "<" + e.GetTagName()

	for key, val := range e.GetAttributes() {
		tagContent = tagContent + " " + key + "=\"" + val + "\""
	}

	tagContent = tagContent + ">" + e.GetContent()

	_, err = out.Write([]byte(tagContent))

	if err != nil {
		return err
	}

	for _, child := range e.GetChildren() {
		ElementToFile(child, out)
	}

	endTagContent := "</" + e.GetTagName() + ">"
	_, err = out.Write([]byte(endTagContent))

	return err
}

/*
	Find the first element matching a given attribute
	Example: form := page.ByAttribute("id","login-form")
*/
func (p *Page) ByAttribute(name, value string) Element {
	return p.root.ByAttribute(name, value)
}

/*
	Find all elements matching a given attribute
	Example: result := page.AllAttribute("class","search-result-div")
*/
func (p *Page) AllByAttribute(name, value string) []Element {
	return p.root.AllByAttribute(name, value)
}

/*
	Find the first element with this id
	Example: form := page.ById("login-form")
*/
func (p *Page) ById(id string) Element {
	return p.root.ById(id)
}

/*
	Find all elements with this id
	Example: result := page.AllById("search-result-div")
*/
func (p *Page) AllById(id string) []Element {
	return p.root.AllById(id)
}

/*
	Find the first element with this class
	Example: form := page.ByClass("container-div")
*/
func (p *Page) ByClass(class string) Element {
	return p.root.ByClass(class)
}

/*
	Find all elements with this class
	Example: result := page.AllById("news-result-div")
*/
func (p *Page) AllByClass(class string) []Element {
	return p.root.AllByClass(class)
}
