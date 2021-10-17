package wayland

import "fmt"

type ErrorHandler interface {
	Handle(err error)
}

type ErrorHandlerFunc func(err error)

func (h ErrorHandlerFunc) Handle(err error) {
	h(err)
}

type PanicOnError struct{}

func (h PanicOnError) Handle(err error) {
	panic(err)
}

type WaylandError struct {
	// ObjectID contains the object where the error occurred.
	ObjectID ObjectID

	// Code contains the error code.
	Code uint32

	// Message contains the error description.
	Message string
}

func (e WaylandError) Error() string {
	return fmt.Sprintf("object %d: %s (code=%08x)", e.ObjectID, e.Message, e.Code)
}
