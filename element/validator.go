package element

import (
	"github.com/tlowry/grawl/util"
	"regexp"
)

// Validators : a generic way to compare html elements in document Searches
type Validator interface {
	Validate(e Element) bool
	GetRegex() *regexp.Regexp
	SetRegex(*regexp.Regexp)
	FirstOnly() bool
	SetFirstOnly(bool)
}

/* Base Validator provides a basic default implementation of the Validator interface*/
type BaseValidator struct {
	Text      string
	regex     *regexp.Regexp
	firstOnly bool
}

/*
	Construct a Validator, text is the search term which can be a literal :
	Example: NewBaseValidator("my-content-div-01")

	or a standard go regular expression:
	Example: NewBaseValidator("my-content-div*")
*/
func NewBaseValidator(text string) (b *BaseValidator, err error) {
	b = &BaseValidator{}
	b.firstOnly = false
	b.Text = text
	if util.ContainsRegex(text) {
		b.regex, err = regexp.Compile(text)
	}

	return b, err
}

func (t *BaseValidator) Validate(e Element) bool {
	return false
}

func (t *BaseValidator) GetRegex() *regexp.Regexp {
	return t.regex
}

func (t *BaseValidator) SetRegex(reg *regexp.Regexp) {
	t.regex = reg
}

func (t *BaseValidator) SetFirstOnly(first bool) {
	t.firstOnly = first
}

func (t *BaseValidator) FirstOnly() bool {
	return t.firstOnly
}

// TagValidator compares based on the value of the tag name
type TagValidator struct {
	*BaseValidator
	wantedTag string
}

/*
	Construct a TagValidator, text is the search term which can be a literal :
	Example: find all menu tags
	NewTagValidator("menu")

	or a standard go regular expression:
	Example: perhaps we want to find all <menu> and <menuitem> tags
	NewTagValidator("menu*")
*/
func NewTagValidator(tagName string) (*TagValidator, error) {
	b, err := NewBaseValidator(tagName)
	v := TagValidator{b, ""}
	v.wantedTag = tagName
	return &v, err
}

func (t TagValidator) Validate(e Element) bool {
	if t.regex != nil {
		if t.regex.MatchString(e.GetTagName()) {
			return true
		}
	} else if t.wantedTag == e.GetTagName() {
		return true
	}
	return false
}

/* Attribute Validator checks if an element contains a matching
attribute value pair*/
type AttributeValidator struct {
	*BaseValidator
	Key, Val string
}

/*
	Construct a AttributeValidator, text is the search term which can be a literal :
	Example: find all items of the "mytextwidget" class
	NewAttributeValidator("mytextwidget")

	or a standard go regular expression:
	Example: perhaps we want to find all elements whose class begins with "mytext"
	NewAttributeValidator("mytext*")
*/
func NewAttributeValidator(key, text string) (*AttributeValidator, error) {
	b, err := NewBaseValidator(text)

	if err == nil {
		a := AttributeValidator{b, key, text}
		return &a, err
	}
	return nil, err

}

func (t AttributeValidator) Validate(e Element) bool {
	// Try to match the regex if present, otherwise assume plain string
	if t.regex != nil {
		if t.regex.MatchString(e.GetAttribute(t.Key)) {
			return true
		}
	} else if e.GetAttribute(t.Key) == t.Val {
		return true
	}
	return false
}
