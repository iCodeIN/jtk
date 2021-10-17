package wayland

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// Ensure Display implements Connection.
var _ Connection = &Display{}

// Display manages a connection to a Wayland display.
type Display struct {
	socket       *net.UnixConn
	display      *WlDisplay
	globals      *Globals
	errorHandler ErrorHandler

	objects      map[ObjectID]Proxy
	objectsMutex sync.RWMutex

	handlers      map[ObjectID][]Handler
	handlersMutex sync.RWMutex

	id uint32
}

// Connect connects to a Wayland display.
func Connect(display string) (*Display, error) {
	socketPath, err := makeSocketPath(display)
	if err != nil {
		return nil, err
	}

	socket, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: socketPath, Net: "unix"})
	if err != nil {
		return nil, err
	}

	if err := socket.SetReadDeadline(time.Time{}); err != nil {
		return nil, err
	}

	wldisplay := &WlDisplay{id: 1}

	objects := make(map[ObjectID]Proxy)
	objects[wldisplay.id] = wldisplay

	handlers := make(map[ObjectID][]Handler)

	conn := &Display{
		socket:       socket,
		display:      wldisplay,
		errorHandler: PanicOnError{},
		objects:      objects,
		handlers:     handlers,
		id:           1,
	}

	globals := &Globals{
		globals: make(map[string]WlRegistryGlobalEvent),
		conn:    conn,
	}

	conn.globals = globals

	return conn, nil
}

func (d *Display) SetErrorHandler(h ErrorHandler) {
	d.errorHandler = h
}

type cbHandler struct {
	ch   chan uint32
	conn *Display
}

func (s *cbHandler) Handle(event Event) {
	switch t := event.(type) {
	case *WlCallbackDoneEvent:
		s.ch <- t.CallbackData
		close(s.ch)
	}
}

// Sync synchronizes the connection.
func (d *Display) Sync() error {
	callback := &WlCallback{id: d.NewID()}
	request := WlDisplaySyncRequest{
		Callback: callback.id,
	}
	d.RegisterProxy(callback)
	handler := &cbHandler{
		ch:   make(chan uint32),
		conn: d,
	}
	d.RegisterHandler(callback.id, handler)
	defer d.UnregisterHandler(callback.id, handler)

	err := d.SendRequest(d.display.id, &request)
	if err != nil {
		return err
	}

	// TODO(jchw): Need to handle connection error case somehow.
	<-handler.ch

	return nil
}

// Globals gets the globals manager.
func (d *Display) Globals() *Globals {
	return d.globals
}

// NewID returns the next ID.
func (d *Display) NewID() ObjectID {
	return ObjectID(atomic.AddUint32(&d.id, 1))
}

// RegisterProxy registers a new proxy.
func (d *Display) RegisterProxy(proxy Proxy) {
	d.objectsMutex.Lock()
	defer d.objectsMutex.Unlock()

	d.objects[proxy.ID()] = proxy
}

// UnregisterProxy unregisters a proxy.
func (d *Display) UnregisterProxy(proxy Proxy) {
	d.UnregisterObject(proxy.ID())
}

// UnregisterObject unregisters an object.
func (d *Display) UnregisterObject(object ObjectID) {
	d.objectsMutex.Lock()
	defer d.objectsMutex.Unlock()

	delete(d.objects, object)
}

// RegisterHandler registers a new event handler.
func (d *Display) RegisterHandler(object ObjectID, handler Handler) {
	d.handlersMutex.Lock()
	defer d.handlersMutex.Unlock()

	d.handlers[object] = append(d.handlers[object], handler)
}

// UnregisterHandler unregisters an event handler.
func (d *Display) UnregisterHandler(object ObjectID, handler Handler) {
	d.handlersMutex.Lock()
	defer d.handlersMutex.Unlock()

	handlers := []Handler{}
	for _, oldhandler := range d.handlers[object] {
		if oldhandler != handler {
			handlers = append(handlers, oldhandler)
		}
	}
	if len(handlers) > 0 {
		d.handlers[object] = handlers
	} else {
		delete(d.handlers, object)
	}
}

// UnregisterHandlers unregisters event handlers for an object.
func (d *Display) UnregisterHandlers(object ObjectID) {
	d.handlersMutex.Lock()
	defer d.handlersMutex.Unlock()

	delete(d.handlers, object)
}

// SendRequest sends a request for a given object.
func (d *Display) SendRequest(id ObjectID, request Request) error {
	return WriteRequest(d.socket, id, request)
}

// Close closes the connection.
func (d *Display) Close() error {
	return d.socket.Close()
}

// PollEvent reads the socket for a new event.
func (d *Display) PollEvent() (ObjectID, Event, error) {
	scanner, err := ReadEvent(d.socket)
	if err != nil {
		return 0, nil, fmt.Errorf("read event: %w", err)
	}

	d.objectsMutex.RLock()
	object, ok := d.objects[ObjectID(scanner.header.ObjectID)]
	d.objectsMutex.RUnlock()

	if !ok {
		return 0, nil, fmt.Errorf("unknown object id: %d", scanner.header.ObjectID)
	}

	event := object.Dispatch(scanner.header.Opcode)
	if event == nil {
		return 0, nil, fmt.Errorf("unknown event opcode %d in event for %d (interface %s)", scanner.header.Opcode, scanner.header.ObjectID, object.Descriptor().Name)
	}

	err = event.Scan(scanner)
	if err != nil {
		return 0, nil, fmt.Errorf("scanning event %s for %d (interface %s): %w", event.MessageName(), scanner.header.ObjectID, object.Descriptor().Name, err)
	}

	return ObjectID(scanner.header.ObjectID), event, nil
}

// DispatchEvent dispatches the provided event.
func (d *Display) DispatchEvent(object ObjectID, event Event) {
	switch t := event.(type) {
	case *WlDisplayDeleteIDEvent:
		d.UnregisterObject(ObjectID(t.ID))
		d.UnregisterHandlers(ObjectID(t.ID))
	case *WlRegistryGlobalEvent:
		d.globals.registerGlobal(t)
	case *WlRegistryGlobalRemoveEvent:
		d.globals.unregisterGlobal(t)
	case *WlDisplayErrorEvent:
		d.errorHandler.Handle(WaylandError{
			ObjectID: t.ObjectID,
			Code:     t.Code,
			Message:  t.Message,
		})
		return
	}

	d.handlersMutex.Lock()
	handlers := make([]Handler, len(d.handlers[object]))
	copy(handlers, d.handlers[object])
	d.handlersMutex.Unlock()

	for _, handler := range handlers {
		handler.Handle(event)
	}
}

// EventLoop runs the Wayland event loop.
func (d *Display) EventLoop() error {
	for {
		object, event, err := d.PollEvent()
		if err != nil {
			// Treat closed socket as success.
			if errors.Is(err, net.ErrClosed) {
				return nil
			}

			return err
		}

		d.DispatchEvent(object, event)
	}
}

func makeSocketPath(display string) (string, error) {
	xdgRuntimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if xdgRuntimeDir == "" {
		return "", errors.New("XDG_RUNTIME_DIR environment variable not set")
	}

	if display == "" {
		display = os.Getenv("WAYLAND_DISPLAY")
	}

	if display == "" {
		display = "wayland-0"
	}

	return xdgRuntimeDir + "/" + display, nil
}
