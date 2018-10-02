// Package errors implements functions to manipulate errors.
package errors

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

// assert Error implements the error interface.
var _ error = &Error{}

type Error struct {
	Message string
	Tags    map[string]string
	Stack   *Stack
}

// New returns an error that formats as the given text.
func New(text string) error {
	return &Error{Message: text, Tags: map[string]string{}}
}

// NewWithTags returns an error that formats as the given text.
func NewWithTags(text string, tags map[string]string) error {
	return &Error{Message: text, Tags: tags}
}

// NewWithTags returns an error that formats as the given text.
func NewNamed(text string, name string) error {
	return &Error{Message: text, Tags: map[string]string{"name": name}}
}

func TagsExtractor(err error) map[string]string {
	if e, ok := err.(*Error); ok {
		return e.Tags
	}
	return nil
}

// Error implements the error interface.
func (e *Error) Error() string {
	b := new(bytes.Buffer)
	e.printStack(b)
	pad(b, ": ")
	b.WriteString(e.Message)
	if b.Len() == 0 {
		return "no error"
	}
	return b.String()
}

// WrapErr returns a Error for the given error and msg.
func WrapErr(err error, msg, name string) error {
	if err == nil {
		return nil
	}
	e := &Error{Message: fmt.Sprintf("%s: %s", msg, err.Error()), Tags: map[string]string{"name": name}}
	e.populateStack()
	return e
}

// NameError returns a Error for the given error and name of error.
func NameError(err error, name string) error {
	if err == nil {
		return nil
	}
	e := &Error{Message: err.Error(), Tags: map[string]string{"name": name}}
	e.populateStack()
	return e
}

// WrapErr returns a Error for the given error and msg.
func Wrap(err error, msg, name string) error {
	if err == nil {
		return nil
	}
	e := &Error{Message: fmt.Sprintf("%s: %s", msg, err.Error())}
	e.populateStack()
	return e
}

// E is a useful func for instantiating Errors.
func E(args ...interface{}) error {
	if len(args) == 0 {
		panic("call to E with no arguments")
	}
	e := &Error{Tags: make(map[string]string)}
	b := new(bytes.Buffer)
	for _, arg := range args {
		switch arg := arg.(type) {
		case string:
			pad(b, ": ")
			b.WriteString(arg)
		case error:
			pad(b, ": ")
			b.WriteString(arg.Error())
		}
	}
	e.Message = b.String()
	e.populateStack()
	return e
}

// populateStack uses the runtime to populate the Error's stack struct with
// information about the current stack.
func (e *Error) populateStack() {
	e.Stack = &Stack{Callers: callers()}
}

// printStack formats and prints the stack for this Error to the given buffer.
// It should be called from the Error's Error method.
func (e *Error) printStack(b *bytes.Buffer) {
	if e.Stack == nil {
		return
	}

	printCallers := callers()

	// Iterate backward through e.Stack.Callers (the last in the stack is the
	// earliest call, such as main) skipping over the PCs that are shared
	// by the error stack and by this function call stack, printing the
	// names of the functions and their file names and line numbers.
	var prev string // the name of the last-seen function
	var diff bool   // do the print and error call stacks differ now?
	for i := 0; i < len(e.Stack.Callers); i++ {
		thisFrame := frame(e.Stack.Callers, i)
		name := thisFrame.Func.Name()

		if !diff && i < len(printCallers) {
			if name == frame(printCallers, i).Func.Name() {
				// both stacks share this PC, skip it.
				continue
			}
			// No match, don't consider printCallers again.
			diff = true
		}

		// Don't print the same function twice.
		// (Can happen when multiple error stacks have been coalesced.)
		if name == prev {
			continue
		}

		// Find the uncommon prefix between this and the previous
		// function name, separating by dots and slashes.
		trim := 0
		for {
			j := strings.IndexAny(name[trim:], "./")
			if j < 0 {
				break
			}
			if !strings.HasPrefix(prev, name[:j+trim]) {
				break
			}
			trim += j + 1 // skip over the separator
		}

		// Do the printing.
		pad(b, separator)
		fmt.Fprintf(b, "%v:%d: ", thisFrame.File, thisFrame.Line)
		if trim > 0 {
			b.WriteString("...")
		}
		b.WriteString(name[trim:])

		prev = name
	}
}

// frame returns the nth frame, with the frame at top of stack being 0.
func frame(callers []uintptr, n int) *runtime.Frame {
	frames := runtime.CallersFrames(callers)
	var f runtime.Frame
	for i := len(callers) - 1; i >= n; i-- {
		var ok bool
		f, ok = frames.Next()
		if !ok {
			break // Should never happen, and this is just debugging.
		}
	}
	return &f
}

// callers is a wrapper for runtime.callers that allocates a slice.
func callers() []uintptr {
	var stk [64]uintptr
	const skip = 4 // Skip 4 stack frames; ok for Error funcs.
	n := runtime.Callers(skip, stk[:])
	return stk[:n]
}

var separator = ":\n\t"

// pad appends str to the buffer if the buffer already has some data.
func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}
