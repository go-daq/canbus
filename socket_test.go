// Copyright 2016 The go-daq Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canbus_test

import (
	"log"
	"reflect"
	"testing"

	"github.com/go-daq/canbus"
)

func TestSocket(t *testing.T) {
	const (
		endpoint = "vcan0"
		N        = 10
		ID       = 128
	)
	var msg = []byte{0, 0xde, 0xad, 0xbe, 0xef}

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
			msg[0] = byte(i)
			_, err := w.Send(ID, msg)
			if err != nil {
				t.Errorf("error send[%d]: %v\n", i, err)
			}
		}
	}()

	for i := 0; i < N; i++ {
		id, data, err := r.Recv()
		if err != nil {
			t.Fatalf("error recv: %v\n", err)
		}
		if id != ID {
			t.Errorf("got id=%v. want=%d\n", id, ID)
		}
		if got, want := int(data[0]), i; got != want {
			t.Errorf("error index: got=%d. want=%d\n", got, want)
		}

		if !reflect.DeepEqual(data[1:], msg[1:]) {
			t.Errorf("error data:\ngot= %v\nwant=%v\n", data[1:], msg[1:])
		}
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	if err := r.Close(); err != nil {
		t.Fatal(err)
	}

}
