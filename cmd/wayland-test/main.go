package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/jchv/jtk/internal/wayland"
)

func shmTempFile(size int64) (*os.File, error) {
	dir := os.Getenv("XDG_RUNTIME_DIR")
	if dir == "" {
		return nil, errors.New("XDG_RUNTIME_DIR is not defined in env")
	}

	file, err := ioutil.TempFile(dir, "wl_shm")
	if err != nil {
		return nil, err
	}

	err = file.Truncate(size)
	if err != nil {
		return nil, err
	}

	err = os.Remove(file.Name())
	if err != nil {
		return nil, err
	}

	return file, nil
}

func drawimg(t int, data []byte) {
	i := 0
	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			data[i+0] = byte(x ^ y ^ t)
			data[i+1] = byte(x/2 ^ y/2 ^ t)
			data[i+2] = byte(x*2 ^ y*2 ^ t)
			data[i+3] = 0xff
			i += 4
		}
	}
}

func main() {
	display := flag.String("display", "", "Wayland socket to connect to, e.g. wayland-0")
	flag.Parse()

	size := 256 * 256 * 4
	file, err := shmTempFile(int64(size))
	if err != nil {
		log.Fatalf("Error creating shared memory tempfile: %v", err)
	}
	defer file.Close()
	data, err := syscall.Mmap(int(file.Fd()), 0, int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		log.Fatalf("Error mapping file into memory: %v", err)
	}
	drawimg(0, data)

	conn, err := wayland.Connect(*display)
	if err != nil {
		log.Fatalf("Error connecting to Wayland compositor: %v", err)
	}

	go func() {
		if err := conn.EventLoop(); err != nil {
			log.Printf("Error in event loop: %v", err)
		}
	}()

	surface, _ := conn.Globals().WlCompositor().CreateSurface(conn)
	pool, _ := conn.Globals().WlShm().CreatePool(conn, wayland.FD(file.Fd()), int32(size))
	buf, _ := pool.CreateBuffer(conn, 0, 256, 256, 256*4, uint32(wayland.WlShmFormatArgb8888))
	xdgsurface, _ := conn.Globals().XdgWmBase().GetXdgSurface(conn, surface.ID())
	toplevel, _ := xdgsurface.GetToplevel(conn)
	toplevel.SetTitle(conn, "Test!")
	toplevel.SetAppID(conn, "wayland-test")
	xdgsurface.SetWindowGeometry(conn, 0, 0, 256, 256)
	surface.Commit(conn)
	surface.Attach(conn, buf.ID(), 0, 0)

	conn.RegisterHandler(xdgsurface.ID(), wayland.HandlerFunc(func(event wayland.Event) {
		switch t := event.(type) {
		case *wayland.XdgSurfaceConfigureEvent:
			xdgsurface.AckConfigure(conn, t.Serial)
			surface.Commit(conn)
		}
	}))

	conn.Sync()

	for i := 0; i < 256; i++ {
		drawimg(i, data)
		surface.Attach(conn, buf.ID(), 0, 0)
		surface.DamageBuffer(conn, 0, 0, 256, 256)
		surface.Commit(conn)
		time.Sleep(time.Second / 30)
	}

	defer conn.Close()
}
