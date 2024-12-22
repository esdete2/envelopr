package handler

import "fmt"

type Error struct {
	Type    ErrorType
	Doc     string
	Wrapped error
}

type ErrorType int

const (
	ErrorLoadingFiles ErrorType = iota
	ErrorRendering
	ErrorCompiling
	ErrorSaving
)

func (e *Error) Error() string {
	switch e.Type {
	case ErrorLoadingFiles:
		return fmt.Sprintf("error loading files: %v", e.Wrapped)
	case ErrorRendering:
		return fmt.Sprintf("error rendering document '%s': %v", e.Doc, e.Wrapped)
	case ErrorCompiling:
		return fmt.Sprintf("error compiling document '%s': %v", e.Doc, e.Wrapped)
	case ErrorSaving:
		return fmt.Sprintf("error saving document '%s': %v", e.Doc, e.Wrapped)
	default:
		return fmt.Sprintf("unknown error: %v", e.Wrapped)
	}
}
