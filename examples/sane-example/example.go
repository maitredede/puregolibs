//go:build linux
// +build linux

// Copyright (C) 2013 Tiago Quelhas. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test program for sane package
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/maitredede/puregolibs/sane"
)

var (
	localOnly bool
)

func main() {
	flag.BoolVar(&localOnly, "localonly", true, "use only local devices")
	flag.Parse()

	if err := sane.Init(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer sane.Exit()

	fmt.Println("looking for devices...")
	devices, err := sane.GetDevices(localOnly)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(devices) == 0 {
		fmt.Println("no devices found")
		os.Exit(2)
	}

	fmt.Printf("found %d devices:\n", len(devices))
	for _, d := range devices {
		fmt.Printf("%s (%s) %s %s\n", d.Name, d.Type, d.Vendor, d.Model)
	}

	first := devices[0]
	h, err := sane.Open(first.Name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer sane.Close(h)

	fmt.Println("starting scan...")
	if err := sane.Start(h); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p, err := sane.GetParameters(h)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("parameters: %+v\n", p)

	bin := &bytes.Buffer{}
	buff := make([]byte, 1024)
	for {
		n, err := sane.Read(h, buff)
		if err == nil {
			data := buff[:n]
			if _, err := bin.Write(data); err != nil {
				panic(err)
			}
			continue
		}
		if err == io.EOF {
			if n != 0 {
				data := buff[:n]
				if _, err := bin.Write(data); err != nil {
					panic(err)
				}
			}
			break
		}
	}

	fmt.Printf("read %d bytes\n", bin.Len())
}
