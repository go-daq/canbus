// Copyright 2016 The go-daq Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canbus_test

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/go-daq/canbus"
)

func TestSocket(t *testing.T) {
	for _, kind := range []canbus.Kind{canbus.SFF, canbus.EFF, canbus.RTR} {
		t.Run(fmt.Sprintf("kind=%v", kind), func(t *testing.T) {
			const (
				endpoint = "vcan0"
				N        = 10
				ID       = 128
			)
			msg := canbus.Frame{
				ID:   ID,
				Data: []byte{0, 0xde, 0xad, 0xbe, 0xef, 0xca, 0xfe, 0xda},
				Kind: kind,
			}
			if kind == canbus.SFF {
				msg.Data = msg.Data[:4]
			}

			r, err := canbus.New()
			if err != nil {
				t.Fatal(err)
			}
			defer r.Close()

			w, err := canbus.New()
			if err != nil {
				t.Fatal(err)
			}
			defer w.Close()

			err = r.Bind(endpoint)
			if err != nil {
				log.Fatal(err)
			}

			err = w.Bind(endpoint)
			if err != nil {
				log.Fatal(err)
			}

			go func() {
				for i := 0; i < N; i++ {
					msg.Data[0] = byte(i)
					_, err := w.Send(msg)
					if err != nil {
						t.Errorf("error send[%d]: %v\n", i, err)
					}
				}
			}()

			for i := 0; i < N; i++ {
				got, err := r.Recv()
				if err != nil {
					t.Fatalf("error recv: %v\n", err)
				}
				if got.ID != ID {
					t.Errorf("got id=%v. want=%d\n", got.ID, ID)
				}
				if got, want := int(got.Data[0]), i; got != want {
					t.Errorf("error index: got=%d. want=%d\n", got, want)
				}
				if got, want := got.Kind, msg.Kind; got != want {
					t.Errorf("error kind: got=%v. want=%v\n", got, want)
				}

				if !reflect.DeepEqual(got.Data[1:], msg.Data[1:]) {
					t.Errorf("error data:\ngot= %v\nwant=%v\n", got.Data[1:], msg.Data[1:])
				}
			}

			if err := w.Close(); err != nil {
				t.Fatal(err)
			}

			if err := r.Close(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestName(t *testing.T) {
	c, err := canbus.New()
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	if got, want := c.Name(), "N/A"; got != want {
		t.Fatalf("invalid socket name: got=%q, want=%q", got, want)
	}

	err = c.Bind("vcan0")
	if err != nil {
		t.Fatalf("could not bind to endpoint: %+v", err)
	}

	if got, want := c.Name(), "vcan0"; got != want {
		t.Fatalf("invalid socket name: got=%q, want=%q", got, want)
	}
}

func TestErrBind(t *testing.T) {
	c, err := canbus.New()
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	err = c.Bind("vcan0_not_there")
	if err == nil {
		t.Fatalf("expected an error binding to invalid CAN bus endpoint")
	}
}

func TestErrLength(t *testing.T) {
	c, err := canbus.New()
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	err = c.Bind("vcan0")
	if err != nil {
		t.Fatalf("could not bind to endpoint: %+v", err)
	}

	_, err = c.Send(canbus.Frame{ID: 42, Data: make([]byte, 8+1)})
	switch {
	case err == nil:
		t.Fatalf("expected an error")
	default:
		if got, want := err.Error(), "canbus: data too big"; got != want {
			t.Fatalf("invalid error: got=%q, want=%q", got, want)
		}
	}
}
