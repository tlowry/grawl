package element

import (
	"code.google.com/p/go.net/html"
	"io"
)

type ElemKind int

const (
	ELEM_VOID ElemKind = 1 + iota
	ELEM_RAW
	ELEM_ESC_RAW
	ELEM_FOREIGN
	ELEM_NORMAL
)

//TODO this is for html5 only
var kindMap map[string]ElemKind

func init() {
	kindMap = map[string]ElemKind{
		"area":     ELEM_VOID,
		"base":     ELEM_VOID,
		"br":       ELEM_VOID,
		"col":      ELEM_VOID,
		"embed":    ELEM_VOID,
		"hr":       ELEM_VOID,
		"img":      ELEM_VOID,
		"input":    ELEM_VOID,
		"keygen":   ELEM_VOID,
		"link":     ELEM_VOID,
		"menuitem": ELEM_VOID,
		"meta":     ELEM_VOID,
		"param":    ELEM_VOID,
		"source":   ELEM_VOID,
		"track":    ELEM_VOID,
		"wbr":      ELEM_VOID,

		"":       ELEM_RAW,
		"script": ELEM_RAW,
		"style":  ELEM_RAW,

		"textarea": ELEM_ESC_RAW,
		"title":    ELEM_ESC_RAW,

		//TODO SVG

		//TODO MATHML
	}

}

func GetElemKind(tagName string) ElemKind {
	kind := kindMap[tagName]
	if kind < 1 {
		return ELEM_NORMAL
	}
	return kind

}

type Parser struct {
	page          *Page
	currentParent Element
	lastElement   Element
}

func NewParser() *Parser {
	p := Parser{}
	return &p
}

// Create a page from a http body
func (p *Parser) ParsePage(r io.Reader) *Page {

	p.page = NewPage()

	tokenizer := html.NewTokenizer(r)
	tokenizer.AllowCDATA(true)
	stillParsing := true

	p.currentParent = p.page.root

	var lastTokenType html.TokenType
	for stillParsing {

		// token type
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			// usually just EOF'
			stillParsing = false

			// carry on if possible TODO Need error here?
		} else {
			token := tokenizer.Token()

			switch tokenType {

			case html.StartTagToken:

				// <tag>
				// type Token struct {
				//     Type     TokenType
				//     DataAtom atom.Atom
				//     Data     string
				//     Attr     []Attribute
				// }
				//
				// type Attribute struct {
				//     Namespace, Key, Val string
				// }
				p.lastElement = p.handleElement(token)

				// Only set this node to current if it can accept child nodes
				if p.lastElement.GetKind() == ELEM_NORMAL {
					p.currentParent = p.lastElement
				}

			case html.TextToken: // text between start and end tag

				if lastTokenType == html.StartTagToken {
					// Last parsed element can support textual content
					p.lastElement.SetContent(token.Data)
				} else if p.currentParent != nil {
					// The last parsed element can't hold text, assume this is for the parent node
					p.currentParent.SetContent(token.Data)
				}
			case html.EndTagToken: // </tag>
				//p.handleElement(token)

				// If a tag has ended it's parent now becomes the parent node
				if p.currentParent.GetParent() != nil {
					p.currentParent = p.currentParent.GetParent()
				}
			case html.SelfClosingTagToken: // <tag/>
				p.lastElement = p.handleElement(token)
			}

		}
		lastTokenType = tokenType

	}

	return p.page

}

func (p *Parser) handleElement(token html.Token) Element {
	foundElem := buildElement(token)

	// Is this the first element?
	if p.page.root == nil {
		if foundElem != nil && foundElem.GetTagName() != "" {
			p.page.root = foundElem
			p.currentParent = foundElem
		}

	} else {
		// Not the first element, it is a child of whatever came before
		p.currentParent.AddChild(foundElem)
	}

	return foundElem
}

// Build an appropriate HTML element from a token
func buildElement(token html.Token) Element {

	var newElem Element

	elemName := token.DataAtom.String()
	switch elemName {
	case "form":
		newElem = NewForm()
	case "input":
		newElem = NewInput()
	default:
		newElem = NewBaseElement()
	}

	newElem.SetTagName(elemName)
	newElem.SetKind(GetElemKind(elemName))

	for _, attr := range token.Attr {
		newElem.SetAttribute(attr.Key, attr.Val)
	}

	return newElem
}
