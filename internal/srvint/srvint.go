package srvint

// #cgo CFLAGS: -I${SRCDIR}/third_party/include
// #cgo LDFLAGS: ${SRCDIR}/third_party/lib/libsrvint.a
// #include "bridge.h"
// #include <stdlib.h>
import "C"

import (
	"context"
	"fmt"
	"unsafe"
)

const maxFrameLength = 6 + 255 + 1

// Server owns a connected libsrvint server transport.
type Server struct {
	ctx *C.srvint_t
}

// New creates a server for the supplied serial device.
func New(device string) (*Server, error) {
	cDevice := C.CString(device)
	defer C.free(unsafe.Pointer(cDevice))
	ctx := C.mbrd_srvint_new(cDevice)
	if ctx == nil {
		return nil, fmt.Errorf("create SrvInt transport: %w", lastError())
	}
	return &Server{ctx: ctx}, nil
}

// Connect opens the configured serial device.
func (s *Server) Connect() error {
	if C.mbrd_srvint_connect(s.ctx) != 0 {
		return fmt.Errorf("connect SrvInt transport: %w", lastError())
	}
	return nil
}

// Serve processes requests until the context is cancelled or transport fails.
func (s *Server) Serve(ctx context.Context) error {
	request := make([]byte, maxFrameLength)
	for {
		if err := ctx.Err(); err != nil {
			return nil
		}
		length := C.mbrd_srvint_receive(s.ctx, (*C.uint8_t)(unsafe.Pointer(&request[0])), C.size_t(len(request)))
		if length < 0 {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("receive SrvInt request: %w", lastError())
		}
		if C.mbrd_srvint_reply(s.ctx, (*C.uint8_t)(unsafe.Pointer(&request[0])), C.size_t(length)) != 0 {
			return fmt.Errorf("reply to SrvInt request: %w", lastError())
		}
	}
}

// Close closes the serial transport. It is safe to call more than once.
func (s *Server) Close() error {
	if s == nil || s.ctx == nil {
		return nil
	}
	result := C.mbrd_srvint_close(s.ctx)
	C.mbrd_srvint_free(s.ctx)
	s.ctx = nil
	if result != 0 {
		return fmt.Errorf("close SrvInt transport: %w", lastError())
	}
	return nil
}

func lastError() error {
	number := C.mbrd_srvint_errno()
	message := C.GoString(C.mbrd_srvint_strerror(number))
	if message == "" {
		message = "unknown error"
	}
	return fmt.Errorf("%s (errno %d)", message, int(number))
}
