// Copyright 2016 The go-daq Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// can-send sends data on the CAN bus.
//
// Usage of can-send:
//
//   can-send [options] <CAN interface> <CAN frame>
//
// where <CAN frame> is of the form: <ID-hex>#<frame data-hex>.
//
// Examples:
//
//  can-send vcan0 f12#1122334455667788
//  can-send vcan0 ffa#deadbeef
//
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-daq/canbus"
)

const (
	frameSep  = "#"
	maxUint32 = uint64(^uint32(0))
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`can-send sends data on the CAN bus.

Usage of can-send:

sh> can-send [options] <CAN interface> <CAN frame>

where <CAN frame> is of the form: <ID-hex>%[1]s<frame data-hex>.

Examples:

 can-send vcan0 f12%[1]s1122334455667788
 can-send vcan0 ffa%[1]sdeadbeef
 `,
			frameSep,
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	log.SetFlags(0)
	log.SetPrefix("can-send> ")

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(2)
	}

	if strings.Count(flag.Arg(1), frameSep) != 1 {
		flag.Usage()
		log.Fatalf(
			"invalid CAN frame (missing %v): %q\n",
			frameSep, flag.Arg(1),
		)
	}

	dev := flag.Arg(0)

	cmd := strings.Split(flag.Arg(1), frameSep)

	if len(cmd[0]) != 3 {
		log.Fatalf("invalid CAN frame id (len=%d != 3)", len(cmd[0]))
	}

	id, err := strconv.ParseUint(cmd[0], 16, 32)
	if err != nil {
		log.Fatalf("error parsing frame id: %v\n", err)
	}

	if id > maxUint32 {
		log.Fatalf("invalid CAN frame id (uint32 overflow): %v\n", id)
	}

	if n := len(cmd[1]); n%2 != 0 {
		log.Fatalf(
			"invalid CAN frame (odd number of bytes): %q (len=%d)\n",
			cmd[1], n,
		)
	}

	data := parseFrame(cmd[1])
	if n, max := len(data), 8; n > max {
		log.Fatalf("invalid CAN frame (len=%d>%d)", n, max)
	}

	sck, err := canbus.New()
	if err != nil {
		log.Fatalf("error creating CAN bus socket: %v\n", err)
	}
	defer sck.Close()

	err = sck.Bind(dev)
	if err != nil {
		log.Fatalf("error binding CAN bus socket: %v\n", err)
	}

	_, err = sck.Send(uint32(id), data)
	if err != nil {
		log.Fatalf("error sending data: %v\n", err)
	}

	err = sck.Close()
	if err != nil {
		log.Fatalf("error closing CAN bus socket: %v\n", err)
	}

}

func parseFrame(str string) []byte {
	const n = 2
	data := make([]byte, 0, len(str)/2)
	for i := 0; i < len(str); i += n {
		v, err := strconv.ParseUint(str[i:i+n], 16, 8)
		if err != nil {
			log.Fatalf("error parsing %q (str[%d:%d]): %v\n",
				str[i:i+n], i, i+n, err,
			)
		}
		data = append(data, byte(v))
	}
	return data
}
