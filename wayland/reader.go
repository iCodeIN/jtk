package wayland

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"syscall"
	"unsafe"
)

var (
	ErrShortRead            = errors.New("short read")
	ErrOutOfBandBufferShort = errors.New("oob buffer short")
	ErrNoOutOfBand          = errors.New("no oob control message")
)

type EventHeader struct {
	ObjectID uint32
	Opcode   uint16
	Size     uint16
}

type EventScanner struct {
	header  EventHeader
	control []syscall.SocketControlMessage
	reader  io.Reader
}

func ReadEvent(socket *net.UnixConn) (*EventScanner, error) {
	buf := [8]byte{}
	oob := [24]byte{}

	n, oobn, _, _, err := socket.ReadMsgUnix(buf[:], oob[:])
	if err != nil {
		return nil, err
	} else if n != 8 {
		return nil, ErrShortRead
	}

	control := []syscall.SocketControlMessage{}
	if oobn > len(oob) {
		return nil, ErrOutOfBandBufferShort
	} else if oobn > 0 {
		control, err = syscall.ParseSocketControlMessage(oob[:])
		if err != nil {
			return nil, fmt.Errorf("parsing socket control message: %w", err)
		}
	}

	header := *(*EventHeader)((unsafe.Pointer)(&buf[0]))
	remaining := int64(header.Size) - int64(len(buf))

	reader := bufio.NewReader(io.LimitReader(socket, remaining))

	return &EventScanner{
		header:  header,
		control: control,
		reader:  reader,
	}, nil
}

func (e *EventScanner) Int() (int32, error) {
	buf := [4]byte{}
	if _, err := e.reader.Read(buf[:]); err != nil {
		return 0, err
	}
	return *(*int32)(unsafe.Pointer(&buf[0])), nil
}

func (e *EventScanner) Uint() (uint32, error) {
	buf := [4]byte{}
	if _, err := e.reader.Read(buf[:]); err != nil {
		return 0, err
	}
	return *(*uint32)(unsafe.Pointer(&buf[0])), nil
}

func (e *EventScanner) ObjectID() (ObjectID, error) {
	buf := [4]byte{}
	if _, err := e.reader.Read(buf[:]); err != nil {
		return 0, err
	}
	return *(*ObjectID)(unsafe.Pointer(&buf[0])), nil
}

func (e *EventScanner) Fixed() (Fixed, error) {
	buf := [4]byte{}
	if _, err := e.reader.Read(buf[:]); err != nil {
		return 0, err
	}
	return *(*Fixed)(unsafe.Pointer(&buf[0])), nil
}

func (e *EventScanner) String() (string, error) {
	len, err := e.Uint()
	if err != nil {
		return "", err
	}

	buf := make([]byte, len+(4-(len&3)))
	if _, err := e.reader.Read(buf[:]); err != nil {
		return "", err
	}

	return string(bytes.TrimRight(buf[:len], "\x00")), nil
}

func (e *EventScanner) Array() ([]byte, error) {
	len, err := e.Uint()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, len)
	if _, err := e.reader.Read(buf[:]); err != nil {
		return nil, err
	}

	return buf, nil
}

func (e *EventScanner) FD() (FD, error) {
	var control syscall.SocketControlMessage

	if len(e.control) < 1 {
		return 0, ErrNoOutOfBand
	}

	control, e.control = e.control[0], e.control[1:]

	fds, err := syscall.ParseUnixRights(&control)
	if err != nil {
		return 0, err
	}

	return FD(fds[0]), nil
}
