// Copyright 2016 The go-daq Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canbus

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"

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
	iface *net.Interface
	addr  *unix.SockaddrCAN
	dev   device
}

// Name returns the device name the socket is bound to.
func (sck *Socket) Name() string {
	if sck.iface == nil {
		return "N/A"
	}
	return sck.iface.Name
}

// SetFilters sets the CAN_RAW_FILTER option and applies the provided
// filters to the underlying socket.
func (sck *Socket) SetFilters(filters []unix.CanFilter) error {
	err := unix.SetsockoptCanRawFilter(sck.dev.fd, unix.SOL_CAN_RAW, unix.CAN_RAW_FILTER, filters)
	if err != nil {
		return fmt.Errorf("could not set CAN filters: %w", err)
	}

	return nil
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

// Send sends the provided frame on the CAN bus.
func (sck *Socket) Send(msg Frame) (int, error) {
	if len(msg.Data) > 8 {
		return 0, errDataTooBig
	}

	switch msg.Kind {
	case SFF:
		msg.ID &= unix.CAN_SFF_MASK
	case EFF:
		msg.ID &= unix.CAN_EFF_MASK
		msg.ID |= unix.CAN_EFF_FLAG
	case RTR:
		msg.ID &= unix.CAN_EFF_MASK
		msg.ID |= unix.CAN_RTR_FLAG
	case ERR:
		msg.ID &= unix.CAN_ERR_MASK
		msg.ID |= unix.CAN_ERR_FLAG
	}

	var frame [frameSize]byte
	binary.LittleEndian.PutUint32(frame[:4], msg.ID)
	frame[4] = byte(len(msg.Data))
	copy(frame[8:], msg.Data)

	return sck.dev.Write(frame[:])
}

// Recv receives data from the CAN socket.
func (sck *Socket) Recv() (msg Frame, err error) {
	var frame [frameSize]byte
	n, err := io.ReadFull(sck.dev, frame[:])
	if err != nil {
		return msg, err
	}

	if n != len(frame) {
		return msg, io.ErrUnexpectedEOF
	}

	msg.ID = binary.LittleEndian.Uint32(frame[:4])
	switch {
	case msg.ID&unix.CAN_EFF_FLAG != 0:
		msg.Kind = EFF
		msg.ID &= unix.CAN_EFF_MASK
	case msg.ID&unix.CAN_ERR_FLAG != 0:
		msg.Kind = ERR
		msg.ID &= unix.CAN_ERR_MASK
	case msg.ID&unix.CAN_RTR_FLAG != 0:
		msg.Kind = RTR
		// FIXME(sbinet): are we sure using the EFF mask makes sense ?
		msg.ID &= unix.CAN_EFF_MASK
	default:
		msg.Kind = SFF
		msg.ID &= unix.CAN_SFF_MASK
	}

	msg.Data = make([]byte, frame[4])
	copy(msg.Data, frame[8:])
	return msg, nil
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
