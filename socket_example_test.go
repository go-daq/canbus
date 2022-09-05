// Copyright 2022 The go-daq Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canbus_test

import (
	"fmt"
	"log"

	"github.com/go-daq/canbus"
	"golang.org/x/sys/unix"
)

func ExampleSocket() {
	recv, err := canbus.New()
	if err != nil {
		log.Fatal(err)
	}
	defer recv.Close()

	send, err := canbus.New()
	if err != nil {
		log.Fatal(err)
	}
	defer send.Close()

	err = recv.Bind("vcan0")
	if err != nil {
		log.Fatalf("could not bind recv socket: %+v", err)
	}

	err = send.Bind("vcan0")
	if err != nil {
		log.Fatalf("could not bind send socket: %+v", err)
	}

	for i := 0; i < 5; i++ {
		_, err := send.Send(canbus.Frame{
			ID:   0x123,
			Data: []byte(fmt.Sprintf("data-%02d", i)),
			Kind: canbus.SFF,
		})
		if err != nil {
			log.Fatalf("could not send frame %d: %+v", i, err)
		}
	}

	for i := 0; i < 5; i++ {
		frame, err := recv.Recv()
		if err != nil {
			log.Fatalf("could not recv frame %d: %+v", i, err)
		}
		fmt.Printf("frame-%02d: %q (id=0x%x)\n", i, frame.Data, frame.ID)
	}

	// Output:
	// frame-00: "data-00" (id=0x123)
	// frame-01: "data-01" (id=0x123)
	// frame-02: "data-02" (id=0x123)
	// frame-03: "data-03" (id=0x123)
	// frame-04: "data-04" (id=0x123)
}

func ExampleSocket_eFF() {
	recv, err := canbus.New()
	if err != nil {
		log.Fatal(err)
	}
	defer recv.Close()

	send, err := canbus.New()
	if err != nil {
		log.Fatal(err)
	}
	defer send.Close()

	err = recv.Bind("vcan0")
	if err != nil {
		log.Fatalf("could not bind recv socket: %+v", err)
	}

	err = send.Bind("vcan0")
	if err != nil {
		log.Fatalf("could not bind send socket: %+v", err)
	}

	for i := 0; i < 5; i++ {
		_, err := send.Send(canbus.Frame{
			ID:   0x1234567,
			Data: []byte(fmt.Sprintf("data-%02d", i)),
			Kind: canbus.EFF,
		})
		if err != nil {
			log.Fatalf("could not send frame %d: %+v", i, err)
		}
	}

	for i := 0; i < 5; i++ {
		frame, err := recv.Recv()
		if err != nil {
			log.Fatalf("could not recv frame %d: %+v", i, err)
		}
		fmt.Printf("frame-%02d: %q (id=0x%x)\n", i, frame.Data, frame.ID)
	}

	// Output:
	// frame-00: "data-00" (id=0x1234567)
	// frame-01: "data-01" (id=0x1234567)
	// frame-02: "data-02" (id=0x1234567)
	// frame-03: "data-03" (id=0x1234567)
	// frame-04: "data-04" (id=0x1234567)
}

func ExampleSocket_setFilters() {
	recv1, err := canbus.New()
	if err != nil {
		log.Fatal(err)
	}
	defer recv1.Close()

	err = recv1.SetFilters([]unix.CanFilter{
		{Id: 0x123, Mask: unix.CAN_SFF_MASK},
	})
	if err != nil {
		log.Fatalf("could not set CAN filters: %+v", err)
	}

	recv2, err := canbus.New()
	if err != nil {
		log.Fatal(err)
	}
	defer recv2.Close()

	err = recv2.SetFilters([]unix.CanFilter{
		{Id: 0x123 | unix.CAN_INV_FILTER, Mask: unix.CAN_SFF_MASK},
	})
	if err != nil {
		log.Fatalf("could not set CAN filters: %+v", err)
	}

	send, err := canbus.New()
	if err != nil {
		log.Fatal(err)
	}
	defer send.Close()

	err = recv1.Bind("vcan0")
	if err != nil {
		log.Fatalf("could not bind recv socket: %+v", err)
	}

	err = recv2.Bind("vcan0")
	if err != nil {
		log.Fatalf("could not bind recv socket: %+v", err)
	}

	err = send.Bind("vcan0")
	if err != nil {
		log.Fatalf("could not bind send socket: %+v", err)
	}

	for i := 0; i < 6; i++ {
		var id uint32 = 0x123
		if i%2 == 0 {
			id = 0x321
		}
		_, err := send.Send(canbus.Frame{
			ID:   id,
			Data: []byte(fmt.Sprintf("data-%02d", i)),
			Kind: canbus.SFF,
		})
		if err != nil {
			log.Fatalf("could not send frame %d: %+v", i, err)
		}
	}

	fmt.Printf("recv1:\n")
	for i := 0; i < 3; i++ {
		frame, err := recv1.Recv()
		if err != nil {
			log.Fatalf("could not recv frame %d: %+v", i, err)
		}
		fmt.Printf("frame-%02d: %q (id=0x%x)\n", i, frame.Data, frame.ID)
	}

	fmt.Printf("\n")
	fmt.Printf("recv2:\n")
	for i := 0; i < 3; i++ {
		frame, err := recv2.Recv()
		if err != nil {
			log.Fatalf("could not recv frame %d: %+v", i, err)
		}
		fmt.Printf("frame-%02d: %q (id=0x%x)\n", i, frame.Data, frame.ID)
	}

	// Output:
	// recv1:
	// frame-00: "data-01" (id=0x123)
	// frame-01: "data-03" (id=0x123)
	// frame-02: "data-05" (id=0x123)
	//
	// recv2:
	// frame-00: "data-00" (id=0x321)
	// frame-01: "data-02" (id=0x321)
	// frame-02: "data-04" (id=0x321)
}
