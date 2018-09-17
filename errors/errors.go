// Package errors implements functions to manipulate errors.
package errors

// New returns an error that formats as the given text.
func New(text string) error {
	return &TaggedError{s: text, tags: map[string]string{}}
}

// NewWithTags returns an error that formats as the given text.
func NewWithTags(text string, tags map[string]string) error {
	return &TaggedError{s: text, tags: tags}
}

// NewWithTags returns an error that formats as the given text.
func NewNamed(text string, name string) error {
	return &TaggedError{s: text, tags: map[string]string{"name": name}}
}

// errorString is a trivial implementation of error.
type TaggedError struct {
	s    string
	tags map[string]string
}

func (e *TaggedError) Error() string {
	return e.s
}

func (e TaggedError) Tags() map[string]string {
	return e.tags
}

func TagsExtractor(err error) map[string]string {
	if e, ok := err.(*TaggedError); ok {
		return e.Tags()
	}
	return nil
}
