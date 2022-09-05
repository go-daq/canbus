// Copyright 2016 The go-daq Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canbus_test

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/go-daq/canbus"
	"golang.org/x/sys/unix"
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

func TestSetFilters(t *testing.T) {
	const (
		endpoint = "vcan0"
		kind     = canbus.EFF
	)

	r1, err := canbus.New()
	if err != nil {
		t.Fatal(err)
	}
	defer r1.Close()

	err = r1.SetFilters([]unix.CanFilter{
		{Id: 0x123, Mask: unix.CAN_SFF_MASK},
	})
	if err != nil {
		t.Fatalf("could not set CAN filters: %+v", err)
	}

	r2, err := canbus.New()
	if err != nil {
		t.Fatal(err)
	}
	defer r2.Close()

	err = r2.SetFilters([]unix.CanFilter{
		{Id: 0xddd, Mask: unix.CAN_SFF_MASK},
	})
	if err != nil {
		t.Fatalf("could not set CAN filters: %+v", err)
	}

	r3, err := canbus.New()
	if err != nil {
		t.Fatal(err)
	}
	defer r3.Close()

	err = r3.SetFilters([]unix.CanFilter{
		{Id: 0xddd | unix.CAN_INV_FILTER, Mask: unix.CAN_SFF_MASK},
	})
	if err != nil {
		t.Fatalf("could not set CAN filters: %+v", err)
	}

	r4, err := canbus.New()
	if err != nil {
		t.Fatal(err)
	}
	defer r2.Close()

	w, err := canbus.New()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	err = r1.Bind(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	err = r2.Bind(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	err = r3.Bind(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	err = r4.Bind(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	err = w.Bind(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 4; i++ {
		var id uint32 = 0x123
		if i%2 == 0 {
			id = 0xddd
		}
		_, err = w.Send(canbus.Frame{
			ID:   id,
			Data: []byte(fmt.Sprintf("%02d", i)),
			Kind: kind,
		})
		if err != nil {
			t.Fatalf("could not send frame %d: %+v", i, err)
		}
	}

	recv := func(r *canbus.Socket, n int) ([]string, error) {
		var msgs []string
		for i := 0; i < n; i++ {
			frame, err := r.Recv()
			if err != nil {
				return nil, fmt.Errorf("could retrieve frame %d: %+v", i, err)
			}
			msgs = append(msgs, string(frame.Data))
		}
		return msgs, nil
	}

	got, err := recv(r1, 2)
	if err != nil {
		t.Fatalf("filtered-socket: %+v", err)
	}

	if got, want := strings.Join(got, ","), "01,03"; got != want {
		t.Fatalf("r1 filter failed:\ngot= %q\nwant=%q\n", got, want)
	}

	got, err = recv(r2, 2)
	if err != nil {
		t.Fatalf("filtered-socket: %+v", err)
	}

	if got, want := strings.Join(got, ","), "00,02"; got != want {
		t.Fatalf("r2 filter failed:\ngot= %q\nwant=%q\n", got, want)
	}

	got, err = recv(r3, 2)
	if err != nil {
		t.Fatalf("filtered-socket: %+v", err)
	}

	if got, want := strings.Join(got, ","), "01,03"; got != want {
		t.Fatalf("r3 filter failed:\ngot= %q\nwant=%q\n", got, want)
	}

	got, err = recv(r4, 4)
	if err != nil {
		t.Fatalf("non-filtered socket: %+v", err)
	}

	if got, want := strings.Join(got, ","), "00,01,02,03"; got != want {
		t.Fatalf("r4 nofilter failed:\ngot= %q\nwant=%q\n", got, want)
	}
}

func BenchmarkSend(b *testing.B) {
	for _, bc := range []struct {
		kind canbus.Kind
		data []byte
	}{
		{canbus.SFF, []byte("0123")},
		{canbus.EFF, []byte("01234567")},
	} {
		b.Run(bc.kind.String(), func(b *testing.B) {
			c, err := canbus.New()
			if err != nil {
				b.Fatal(err)
			}
			defer c.Close()

			err = c.Bind("vcan0")
			if err != nil {
				b.Fatal(err)
			}

			frame := canbus.Frame{Kind: bc.kind, Data: bc.data, ID: 0xff}
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, err = c.Send(frame)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkSendRecv(b *testing.B) {
	for _, bc := range []struct {
		kind canbus.Kind
		data []byte
	}{
		{canbus.SFF, []byte("0123")},
		{canbus.EFF, []byte("01234567")},
	} {
		b.Run(bc.kind.String(), func(b *testing.B) {
			w, err := canbus.New()
			if err != nil {
				b.Fatal(err)
			}
			defer w.Close()

			err = w.Bind("vcan0")
			if err != nil {
				b.Fatal(err)
			}

			r, err := canbus.New()
			if err != nil {
				b.Fatal(err)
			}
			defer r.Close()

			err = r.Bind("vcan0")
			if err != nil {
				b.Fatal(err)
			}

			frame := canbus.Frame{Kind: bc.kind, Data: bc.data, ID: 0xff}
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, err = w.Send(frame)
				if err != nil {
					b.Fatal(err)
				}

				_, err = r.Recv()
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
