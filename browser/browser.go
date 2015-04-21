/*
	A grawl browser can load pages and submit forms
	found on those pages.
*/
package browser

import (
	"fmt"
	"github.com/tlowry/grawl/element"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

const (
	UAgent string = "Grawl v0.1"
)

type Browser struct {
	Client    *http.Client
	url       string
	userAgent string
}

// Create a new Grawl Browser
func NewBrowser() *Browser {
	b := Browser{}
	b.Client = &http.Client{}

	var err error
	b.Client.Jar, err = cookiejar.New(nil)

	if err != nil {
		panic(fmt.Sprintf("Failed to create a cookie jar for browser %s ", err.Error()))
	}

	b.userAgent = UAgent

	return &b
}

// Create a grawl browser using a predefined http client
func NewBrowserWithClient(client *http.Client) *Browser {
	b := Browser{}
	b.Client = client

	return &b
}

// Return the user agent the browser is currently using in requests
func (b *Browser) GetUserAgent() string {
	return b.userAgent
}

// Set the user agent to be used in browser requests
func (b *Browser) SetUserAgent(agent string) {
	b.userAgent = agent
}

// Post a form to the site this browser is connected to
func (b *Browser) SubmitForm(form *element.Form) *element.Page {

	method := form.Method()

	if method == "" || len(method) < 1 {
		method = "POST"
	} else {
		method = "GET"
	}

	action := b.RelToAbs(form.Action())
	action = FixProtocol(action)

	// gather the form values
	vals := url.Values{}
	for _, in := range form.GetInputs() {
		if len(in.GetAttribute("name")) > 0 {
			vals.Set(in.GetAttribute("name"), in.GetAttribute("value"))
		}
	}

	for _, in := range form.GetSelects() {
		if len(in.GetAttribute("name")) > 0 {
			vals.Set(in.GetAttribute("name"), in.GetAttribute("value"))
		}
	}

	for x := range vals {
		fmt.Println(x, "=", vals[x])
	}

	// build the request using the method, action and form values
	var resp *http.Response

	var body *strings.Reader = nil

	if method == "POST" {
		body = strings.NewReader(vals.Encode())
	} else {
		action = action + "?" + vals.Encode()
	}

	var req *http.Request
	var err error
	if body == nil {
		req, err = http.NewRequest(strings.ToUpper(method), action, nil)
	} else {
		req, err = http.NewRequest(strings.ToUpper(method), action, body)
	}

	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", b.GetUserAgent())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	log.Println("Sending request " + action)
	resp, err = b.Client.Do(req)
	log.Println("Request sent")
	if err != nil {
		panic(err)
	}

	page := element.ParseResp(resp)
	log.Println("parsed")
	return page

}

// Load a page from a url
func (b *Browser) Load(url string) *element.Page {
	b.url = url
	// Fills in "http://" if the url is missing the protocol
	url = FixProtocol(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(fmt.Sprintf("Error creating request %s", err.Error()))
	}

	req.Header.Set("User-Agent", b.GetUserAgent())

	var resp *http.Response
	resp, err = b.Client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Error submitting request %s", err.Error()))
	}

	page := element.ParseResp(resp)
	page.SetUrl(url)

	return page
}

// Load a page from a local file
func (b *Browser) LoadFile(fileName string) *element.Page {
	b.url = fileName

	file, err := os.Open(fileName)
	if err != nil {
		panic(fmt.Sprintf("Error opening file %s", err.Error()))
	}

	page := element.ParseBody(file)

	return page

}

// Convert a relative url to an absolute url based on the browsers currently loaded page url
func (b *Browser) RelToAbs(relUrl string) string {
	protoAndURL := strings.Split(relUrl, "://")

	fixedUrl := relUrl
	// Is there a protocol string in the url?
	if len(protoAndURL) < 2 {

		// No protocol, remove excess / if required and concatenate
		if strings.HasPrefix(relUrl, "/") && strings.HasSuffix(b.url, "/") {
			fixedUrl = b.url + relUrl[1:len(relUrl)]
		} else {
			fixedUrl = b.url + relUrl
		}

	}
	return fixedUrl
}

func (b *Browser) GetCookies() []*http.Cookie {
	currentURL, _ := url.Parse(b.url)
	return b.Client.Jar.Cookies(currentURL)
}

func (b *Browser) ClearCookies() {
	b.Client.Jar, _ = cookiejar.New(nil)
}

func (b *Browser) SetCookie(cookie *http.Cookie) {
	currentURL, _ := url.Parse(b.url)

	// ToDo lock here in case currentURL changes
	cookies := b.Client.Jar.Cookies(currentURL)
	cookies = append(cookies, cookie)
	b.Client.Jar.SetCookies(currentURL, cookies)
	//ToDo unlock here

}

// prepends http:// to the start of urls which are missing a protocol
func FixProtocol(url string) string {
	protoAndUrl := strings.Split(url, ":")
	if len(protoAndUrl) < 2 {
		url = "http://" + url
	}
	return url
}
