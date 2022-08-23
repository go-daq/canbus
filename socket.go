// Copyright 2016 The go-daq Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package canbus provides high-level access to CAN bus sockets.
//
// A typical usage might look like:
//
//	sck, err := canbus.New()
//	err = sck.Bind("vcan0")
//	for {
//	    id, data, err := sck.Recv()
//	}
package canbus

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"unsafe"

	"golang.org/x/sys/unix"
)

var (
	errDataTooBig = errors.New("canbus: data too big")
)

// New returns a new CAN bus socket.
func New() (*Socket, error) {
	fd, err := unix.Socket(unix.AF_CAN, unix.SOCK_RAW, unix.CAN_RAW)
	if err != nil {
		return nil, err
	}

	return &Socket{dev: device{fd}}, nil
}

// Socket is a high-level representation of a CANBus socket.
type Socket struct {
	iface   *net.Interface
	addr    *unix.SockaddrCAN
	dev     device
	rfilter []unix.CanFilter
}

// Name returns the device name the socket is bound to.
func (sck *Socket) Name() string {
	if sck.iface == nil {
		return "N/A"
	}
	return sck.iface.Name
}

// Close closes the CAN bus socket.
func (sck *Socket) Close() error {
	return unix.Close(sck.dev.fd)
}

// Bind binds the socket on the CAN bus with the given address.
//
// Example:
//
//	err = sck.Bind("vcan0")
func (sck *Socket) Bind(addr string) error {
	iface, err := net.InterfaceByName(addr)
	if err != nil {
		return err
	}

	sck.iface = iface
	sck.addr = &unix.SockaddrCAN{Ifindex: sck.iface.Index}

	return unix.Bind(sck.dev.fd, sck.addr)
}

// ApplyFilters sets the CAN_RAW_FILTER option of the socket and stores the new filters
// in the Socket object for use by other methods such as Recv.
func (sck *Socket) ApplyFilters(filters []unix.CanFilter) {
	unix.SetsockoptCanRawFilter(sck.dev.fd, unix.SOL_CAN_RAW, unix.CAN_RAW_FILTER, filters)
	sck.rfilter = filters
}

// Send sends data with a CAN_frame id to the CAN bus.
// idMask should be a mask such as unix.CAN_EFF_MASK.
func (sck *Socket) Send(id uint32, idMask uint32, data []byte) (int, error) {
	if len(data) > 8 {
		return 0, errDataTooBig
	}

	id &= idMask
	var frame [frameSize]byte
	binary.LittleEndian.PutUint32(frame[:4], id)
	frame[4] = byte(len(data))
	copy(frame[8:], data)

	return sck.dev.Write(frame[:])
}

// Recv receives data from the CAN socket.
// id is the CAN_frame id the data was originated from.
func (sck *Socket) Recv() (id uint32, data []byte, err error) {
	var frame [frameSize]byte
	n, err := io.ReadFull(sck.dev, frame[:])
	if err != nil {
		return id, data, err
	}

	if n != len(frame) {
		return id, data, io.ErrUnexpectedEOF
	}

	id = binary.LittleEndian.Uint32(frame[:4])
	for _, filter := range sck.rfilter {
		if id&filter.Mask == filter.Id&filter.Mask {
			id &= filter.Mask
			break
		}
	}
	data = make([]byte, frame[4])
	copy(data, frame[8:])
	return id, data, nil
}

type device struct {
	fd int
}

func (d device) Read(data []byte) (int, error) {
	return unix.Read(d.fd, data)
}

func (d device) Write(data []byte) (int, error) {
	return unix.Write(d.fd, data)
}

const frameSize = unsafe.Sizeof(frame{})

// frame is a can_frame.
type frame struct {
	ID   uint32
	Len  byte
	_    [3]byte
	Data [8]byte
}
