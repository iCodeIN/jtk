package wayland

//go:generate go run ../cmd/waygen ../third_party/wayland/protocol ../third_party/wayland-protocols

// FD represents a UNIX file descriptor. This type is present inside Wayland
// requests and events, but it is not sent over the main connection, and as
// such is not encoded/decoded into the wire directly.
type FD uintptr

// Message is an interface implemented by all Wayland messages.
type Message interface {
	Opcode() uint16
}