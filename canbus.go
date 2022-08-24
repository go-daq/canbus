// Copyright 2016 The go-daq Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package canbus provides high-level access to CAN bus sockets.
//
// A typical usage might look like:
//
//	sck, err := canbus.New()
//	err = sck.Bind("vcan0")
//	for {
//	    msg, err := sck.Recv()
//	}
package canbus
