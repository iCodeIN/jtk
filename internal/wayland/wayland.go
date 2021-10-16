package wayland

//go:generate go run ../../cmd/waygen ../../third_party/wayland/protocol ../../third_party/wayland-protocols

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
	MessageName() string
}

// NewProxy is a function that can construct a new proxy with a given object ID.
type NewProxy func(ObjectID) Proxy

// Proxy is an interface implemented for proxying server objects.
type Proxy interface {
	// ID returns the object ID of the proxied object.
	ID() ObjectID

	// Descriptor returns the interface descriptor that corresponds to this
	// proxy.
	Descriptor() *InterfaceDescriptor

	Dispatcher
}

// Dispatcher is an interface for something that can dispatch Wayland events.
type Dispatcher interface {
	// Dispatch returns an event, given an event opcode.
	Dispatch(uint16) Event
}

// Handler is an interface that handles events.
type Handler interface {
	Handle(event Event)
}

// HandlerFunc is a helper for using a function as an event handler.
type HandlerFunc func(event Event)

func (f HandlerFunc) Handle(event Event) { f(event) }

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
	Interfaces []*InterfaceDescriptor
}

// InterfaceDescriptor contains runtime metadata about an interface.
type InterfaceDescriptor struct {
	Name     string
	Events   []EventDescriptor
	Requests []RequestDescriptor
	NewProxy NewProxy
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

// Connection is a type implemented by a Wayland connection manager.
type Connection interface {
	// NewID returns the next object ID.
	NewID() ObjectID

	// RegisterProxy registers a new proxy.
	RegisterProxy(Proxy)

	// UnregisterProxy unregisters a proxy.
	UnregisterProxy(Proxy)

	// SendRequest sends a request for a given object.
	SendRequest(ObjectID, Request) error
}
