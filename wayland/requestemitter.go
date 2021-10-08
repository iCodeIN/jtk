package wayland

import (
	"bytes"
	"errors"
	"io"
	"net"
	"syscall"
	"unsafe"
)

var (
	ErrMessageOverflow = errors.New("message too large")
)

type RequestHeader struct {
	ObjectID uint32
	Opcode   uint16
	Size     uint16
}

type RequestEmitter struct {
	writer io.Writer
	oob    []byte
}

func WriteRequest(socket *net.UnixConn, object ObjectID, request Request) error {
	w := bytes.Buffer{}

	w.Write(make([]byte, 8))

	emitter := RequestEmitter{&w, nil}
	if err := request.Emit(&emitter); err != nil {
		return err
	}

	size := w.Len() + 8
	if size > int(uint16(size)) || size < 8 {
		return ErrMessageOverflow
	}

	buf := w.Bytes()

	*(*RequestHeader)((unsafe.Pointer)(&buf[0])) = RequestHeader{
		ObjectID: uint32(object),
		Opcode:   request.Opcode(),
		Size:     uint16(w.Len() + 8),
	}

	n, oobn, err := socket.WriteMsgUnix(buf, emitter.oob, nil)
	if err != nil {
		return err
	}
	if oobn != len(emitter.oob) || n != len(buf) {
		return io.ErrShortWrite
	}

	return nil
}

func (e *RequestEmitter) PutInt(v int32) error {
	buf := [4]byte{}
	*(*int32)(unsafe.Pointer(&buf[0])) = v
	_, err := e.writer.Write(buf[:])
	return err
}

func (e *RequestEmitter) PutUint(v uint32) error {
	buf := [4]byte{}
	*(*uint32)(unsafe.Pointer(&buf[0])) = v
	_, err := e.writer.Write(buf[:])
	return err
}

func (e *RequestEmitter) PutObjectID(v ObjectID) error {
	buf := [4]byte{}
	*(*ObjectID)(unsafe.Pointer(&buf[0])) = v
	_, err := e.writer.Write(buf[:])
	return err
}

func (e *RequestEmitter) PutFixed(v Fixed) error {
	buf := [4]byte{}
	*(*Fixed)(unsafe.Pointer(&buf[0])) = v
	_, err := e.writer.Write(buf[:])
	return err
}

func (e *RequestEmitter) PutString(v string) error {
	if err := e.PutUint(uint32(len(v))); err != nil {
		return err
	}

	if _, err := e.writer.Write([]byte(v)); err != nil {
		return err
	}

	pad := 4 - (len(v) & 0x3)
	if pad > 0 {
		if _, err := e.writer.Write(make([]byte, pad)); err != nil {
			return err
		}
	}

	return nil
}

func (e *RequestEmitter) PutArray(v []byte) error {
	if err := e.PutUint(uint32(len(v))); err != nil {
		return err
	}
	_, err := e.writer.Write(v)
	return err
}

func (e *RequestEmitter) PutFD(v FD) error {
	e.oob = append(e.oob, syscall.UnixRights(int(v))...)
	return nil
}
