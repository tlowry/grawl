package element

import (
	"github.com/hishboy/gocommons/lang"
	"sync"
)

/*
	Define common methods for HTML elements
*/
type Element interface {
	AddChild(child Element)
	RemoveChild(child Element)
	GetParent() Element
	setParent(newParent Element)
	GetAttribute(key string) string
	SetAttribute(key, value string)
	RemoveAttribute(key string)
	GetTagName() string
	SetTagName(name string)
	GetChildren() []Element
	GetContent() string
	SetContent(string)
	String() string
	GetAttributes() map[string]string
	ByAttribute(name, value string) Element
	AllByAttribute(name, value string) []Element
	ById(id string) Element
	AllById(id string) []Element
	ByClass(class string) Element
	AllByClass(class string) []Element
	ByTag(tag string) Element
	AllByTag(tag string) []Element
	GetMutex() *sync.RWMutex
	GetKind() ElemKind
	SetKind(ElemKind)
	Next() Element
	setNext(Element)
	Prev() Element
	setPrev(Element)
}

/*
	Common base HTML Element implementation
*/
type BaseElement struct {
	attributes map[string]string
	children   []Element
	parent     Element
	tagName    string
	data       string
	mutex      *sync.RWMutex
	kind       ElemKind
	content    string
	next       Element
	prev       Element
}

func NewBaseElement() *BaseElement {
	e := BaseElement{}
	e.tagName = "Unknown"
	e.parent = nil
	e.attributes = make(map[string]string)
	e.children = []Element{}
	e.mutex = &sync.RWMutex{}
	e.kind = ELEM_NORMAL
	return &e
}

// Append a child to a node and remove it from its old parent if it has one
func (e *BaseElement) AddChild(child Element) {

	// Remove child from current parent if needed
	currentParent := child.GetParent()
	if currentParent != nil {
		currentParent.RemoveChild(child)
	}

	// Make this child a sibling of existing children
	children := e.GetChildren()
	if len(children) > 0 {
		lastElem := children[len(children)-1]
		child.setPrev(lastElem)
	}
	e.children = append(e.children, child)
	child.setNext(nil)
	child.setParent(e)
}

// Remove a child and all its children from a given elements tree
func (e *BaseElement) RemoveChild(child Element) {
	for i, c := range e.GetChildren() {
		if child == c {
			// TODO need to lock here
			e.children = append(e.children[:i], e.children[i+1:]...)
			child.setParent(nil)
		}
	}
}

/*
	Change the nodes parent node
	(should be private to prevent users leaving nodes
	under multiple parents child slices)
*/

func (e *BaseElement) setParent(newParent Element) {
	e.parent = newParent
}

// Return this childs current parent element
func (e *BaseElement) GetParent() Element {
	return e.parent
}

/*
	Return an Elements tag attribute
	Example: Read the css style attribute if present
	styleStr := elem.GetAttribute("style")
*/
func (e *BaseElement) GetAttribute(key string) string {
	return e.attributes[key]
}

/*
	Set an Elements tag attribute
	Example: Use the style attribute to hide the element
	elem.SetAttribute("style","display:none")
*/
func (e *BaseElement) SetAttribute(key, value string) {
	e.attributes[key] = value
}

/*
	Get all tag attributes belonging to an element
*/
func (e *BaseElement) GetAttributes() map[string]string {
	return e.attributes
}

/*
	Completely remove an attribute from an element
	Example: Remove all local formatting from an element
	elem.RemoveAttribute("style")
*/
func (e *BaseElement) RemoveAttribute(key string) {
	delete(e.attributes, key)
}

// Return all children directly below this element
func (e *BaseElement) GetChildren() []Element {
	return e.children
}

// Return the tag name for an element such as "img" or "div"
func (e *BaseElement) GetTagName() string {
	return e.tagName
}

// Set the tag name of this Element to the given name
func (e *BaseElement) SetTagName(name string) {
	e.tagName = name
}

// Return a textual representation of this element
func (e *BaseElement) String() string {

	ret := "<" + e.GetTagName()
	for key, val := range e.GetAttributes() {
		ret = ret + " " + key + "=\"" + val + "\""
	}
	ret = ret + "> parent: "
	parent := e.GetParent()
	if parent == nil {
		ret = ret + "nil"

	} else {
		ret = ret + parent.GetTagName()
	}
	return ret
}

/*
	Returns text enclosed within this elements start and end tags
	excluding any child tags.
*/
func (e *BaseElement) GetContent() string {
	return e.content
}

/*
	Sets the text enclosed within this elements start and end tags
	excluding any child tags.
*/
func (e *BaseElement) SetContent(content string) {
	e.content = content
}

/*
	Return the first element found with a given
	attribute matching the given value.
	Usually used when there is only one of these
	elements on a page or any one of them will do.

	Example: Return the first link with a url
	pointing to /search
	link := elem.ByAttribute("href", "/search")
*/
func (e *BaseElement) ByAttribute(name, value string) Element {
	val, err := NewAttributeValidator(name, value)

	if err == nil {
		elem := BFSFirst(e, val)
		return elem
	}

	panic(err)
}

/*
	Return all elements found with the given
	attribute matching the given value.

	Example:
	find all links where the url points to submit
	elem.ByAttribute("href", "/submit")
*/
func (e *BaseElement) AllByAttribute(name, value string) []Element {
	val, err := NewAttributeValidator(name, value)

	if err == nil {
		return BFS(e, val)
	}

	panic(err)

}

/*
	Return the first element found with the given
	id.

	Example:
	Find the first (usually only) search box div
	box := elem.ById("search-box")
*/
func (e *BaseElement) ById(id string) Element {
	val, err := NewAttributeValidator("id", id)
	if err == nil {
		return BFSFirst(e, *val)
	}

	panic(err)

}

/*
	Return the all element found with the given
	id.

	Example:
	find all links where the url points to submit
	elem := elem.ById("search-result")
*/
func (e *BaseElement) AllById(id string) []Element {
	val, err := NewAttributeValidator("id", id)
	if err == nil {
		return BFS(e, val)
	}

	panic(err)
}

/*
	Return the first element found with the given
	id.

	Example:
	Find the first (usually only) search box div
	box := elem.ByClass("info-container")
*/
func (e *BaseElement) ByClass(class string) Element {
	val, err := NewAttributeValidator("class", class)
	if err == nil {
		return BFSFirst(e, *val)
	}

	panic(err)

}

/*
	Return the all element found with the given
	class.

	Example:
	elem := elem.AllByClass("info-result")
*/
func (e *BaseElement) AllByClass(class string) []Element {
	val, err := NewAttributeValidator("class", class)

	if err == nil {
		return BFS(e, val)
	}

	panic(err)
}

/*
	Return the first element with a matching tag
	Example: find the first form in the document
	form := elem.ByTag("form")
*/
func (e *BaseElement) ByTag(tag string) Element {
	val, err := NewTagValidator(tag)

	if err == nil {
		return BFSFirst(e, *val)
	}

	panic(err)
}

/*
	Return the first element with a matching tag
	Example: find all forms in the document
	forms := elem.ByTag("form")
*/
func (e *BaseElement) AllByTag(tag string) []Element {
	val, err := NewTagValidator(tag)

	if err == nil {
		return BFS(e, *val)
	}

	panic(err)
}

func (e *BaseElement) GetMutex() *sync.RWMutex {
	return e.mutex
}

func (e *BaseElement) GetKind() ElemKind {
	return e.kind
}

func (e *BaseElement) SetKind(elemKind ElemKind) {
	e.kind = elemKind
}

func (e *BaseElement) Next() Element {
	return e.next
}

func (e *BaseElement) setNext(next Element) {
	e.next = next
}

func (e *BaseElement) Prev() Element {
	return e.prev
}

func (e *BaseElement) setPrev(next Element) {
	e.prev = next
}

// Perform an iterative Breadth first search to find a matching element
func BFSFirst(current Element, v Validator) Element {
	v.SetFirstOnly(true)
	results := BFS(current, v)
	if len(results) > 0 {
		return results[0]
	}
	return nil
}

/*
	Perform an iterative Breadth first search to find all matching elements.
	Searching begins at the given element and all subchildren of this Element
	will be
*/
func BFS(current Element, v Validator) []Element {
	queue := lang.NewQueue()

	queue.Push(current)

	results := []Element{}

	for current != nil {

		tmp := queue.Poll()

		if tmp != nil {
			current = tmp.(Element)
			//current.GetMutex().RLock()

			if v.Validate(current) {
				results = append(results, current)
				if v.FirstOnly() {
					break
				}
			}
			for _, child := range current.GetChildren() {
				queue.Push(child)
			}

			//current.GetMutex().RUnlock()

		} else {

			current = nil
		}
	}

	return results
}
