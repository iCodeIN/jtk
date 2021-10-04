package wayland

import (
	"errors"
	"net"
	"os"
	"time"
)

type Display struct {
	socket *net.UnixConn
}

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

	return &Display{
		socket: socket,
	}, nil
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

func (d *Display) Close() error {
	return d.socket.Close()
}
