// Copyright 2016 The go-daq Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// can-dump prints data flowing on CAN bus.
//
// Usage of can-dump:
//
//  can-dump [options] <CAN interface>
//    (use CTRL-C to terminate can-dump)
//
// Examples:
//
//   can-dump vcan0
//
package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-daq/canbus"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`can-dump prints data flowing on CAN bus.

Usage of can-dump:

  can-dump [options] <CAN interface>
    (use CTRL-C to terminate can-dump)

Examples:

 can-dump vcan0
`,
		)
		flag.PrintDefaults()
		os.Exit(2)
	}

	flag.Parse()

	log.SetFlags(0)
	log.SetPrefix("can-dump> ")

	if flag.NArg() < 1 {
		flag.Usage()
	}

	sck, err := canbus.New()
	if err != nil {
		log.Fatal(err)
	}
	defer sck.Close()

	addr := flag.Arg(0)
	err = sck.Bind(addr)
	if err != nil {
		log.Fatalf("error binding to [%s]: %v\n", addr, err)
	}

	var blank = strings.Repeat(" ", 24)
	for {
		id, msg, err := sck.Recv()
		if err != nil {
			log.Fatalf("recv error: %v\n", err)
		}
		ascii := strings.ToUpper(hex.Dump(msg))
		ascii = strings.TrimRight(strings.Replace(ascii, blank, "", -1), "\n")
		fmt.Printf("%7s  %03x %s\n", sck.Name(), id, ascii)
	}
}
