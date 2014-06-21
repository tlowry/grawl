package element

import (
	"net/url"
)

// HTML <input> object
type Input struct {
	BaseElement
}

func NewInput() *Input {
	e := Input{*NewBaseElement()}
	return &e
}

/*
	HTML <form> object
*/
type Form struct {
	BaseElement
	Inputs *url.Values
}

func NewForm() *Form {
	form := Form{*NewBaseElement(), &url.Values{}}
	return &form
}

func (e *Form) GetInputs() []Element {
	val, err := NewTagValidator("input")
	if err == nil {
		return BFS(e, *val)
	}

	return []Element{}
}

func (e *Form) Name() string {
	return e.GetAttribute("name")
}

func (e *Form) Method() string {
	return e.GetAttribute("method")
}

func (e *Form) Action() string {
	return e.GetAttribute("action")
}

// Todo consider letting users add fields not present
func (e *Form) SetField(name, value string) {
	for _, in := range e.GetInputs() {
		if in.GetAttribute("name") == name {
			in.SetAttribute("value", value)
			break
		}
	}
}

func (e *Form) GetField(name string) string {

	for _, in := range e.GetInputs() {
		if in.GetAttribute("name") == name {
			return in.GetAttribute("value")
		}
	}
	return ""
}
