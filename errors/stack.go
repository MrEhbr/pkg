package errors

// Stack represents errors stack trace
type Stack struct {
	Callers []uintptr
}
