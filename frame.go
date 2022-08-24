// Copyright 2022 The go-daq Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canbus

//go:generate stringer -output=frame_string.go -type Kind

import "unsafe"

// Frame is exchanged over a CAN bus.
type Frame struct {
	ID   uint32
	Data []byte
	Kind Kind
}

type Kind uint8

const (
	SFF Kind = iota // Standard frame format
	EFF             // Extended frame format
	RTR             // Remote transmission request
	ERR             // Error message frame
)

const frameSize = unsafe.Sizeof(
	// this is a can_frame.
	struct {
		ID   uint32
		Len  byte
		_    [3]byte
		Data [8]byte
	}{},
)
