package wayland

//go:generate go run ../cmd/waygen ../third_party/wayland/protocol ../third_party/wayland-protocols

// FD represents a UNIX file descriptor. This type is present inside Wayland
// requests and events, but it is not sent over the main connection, and as
// such is not encoded/decoded into the wire directly.
type FD uintptr

// Fixed is Wayland's fixed-point decimal type.
type Fixed uint32

// ObjectID is an incrementing, per-connection object ID.
type ObjectID uint32

// Message is an interface implemented by all Wayland messages.
type Message interface {
	Opcode() uint16
}

// Dispatcher is an interface implemented by all Wayland interfaces.
type Dispatcher interface {
	// Dispatch returns an event, given an event opcode.
	Dispatch(uint16) Event
}

// Event is an interface implemented by all Wayland events.
type Event interface {
	Message
	Scan(s *EventScanner) error
}

// Request is an interface implemented by all Wayland requests.
type Request interface {
	Message
	Emit(e *RequestEmitter) error
}

// ProtocolDescriptor contains runtime metadata about a protocol.
type ProtocolDescriptor struct {
	Name       string
	Interfaces []InterfaceDescriptor
}

// InterfaceDescriptor contains runtime metadata about an interface.
type InterfaceDescriptor struct {
	Name       string
	Events     []EventDescriptor
	Requests   []RequestDescriptor
	Dispatcher Dispatcher
}

// EventDescriptor contains runtime metadata about an event.
type EventDescriptor struct {
	Name   string
	Opcode uint32
	Type   Event
}

// RequestDescriptor contains runtime metadata about a request.
type RequestDescriptor struct {
	Name   string
	Opcode uint32
	Type   Request
}
