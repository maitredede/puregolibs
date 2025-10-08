//go:build linux
// +build linux

// Copyright (C) 2013 Tiago Quelhas. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test program for sane package
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/maitredede/puregolibs/sane"
)

var (
	localOnly bool
)

func main() {
	flag.BoolVar(&localOnly, "localonly", true, "use only local devices")
	flag.Parse()
	slog.SetLogLoggerLevel(slog.LevelDebug)

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
	fmt.Printf("using first: %s\n", first.Name)
	h, err := sane.Open(first.Name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer h.Close()

	fmt.Println("options :")
	descList, err := h.GetOptionDescriptors()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var curGroup string
	for _, desc := range descList {
		if desc.Type == sane.TypeGroup {
			curGroup = desc.Title
			continue
		}

		s := fmt.Sprintf("%s\t[%2d]'%s' %s %s %d", curGroup, desc.Number, desc.Name, desc.Title, desc.Type, desc.BinSize)

		value, err := desc.GetValue()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		s = fmt.Sprintf("%s\tv=%v", s, value)
		if desc.Constraint != nil {
			s = fmt.Sprintf("%s\tconstraint %s=%v", s, desc.ConstraintType, desc.Constraint)
		}
		fmt.Println(s)
	}

	if _, err := h.SetOptionValue("mode", "Color"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if _, err := h.SetOptionValue("source", "Flatbed"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if _, err := h.SetOptionValue("preview", false); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("starting scan...")
	img, err := h.ReadImage()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("scanned image %v\n", img.Bounds())

	// fmt.Println("starting scan...")
	// if err := sane.Start(h); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// p, err := sane.GetParameters(h)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("parameters: %+v\n", p)

	// bin := &bytes.Buffer{}
	// buff := make([]byte, 1024)
	// for {
	// 	n, err := sane.Read(h, buff)
	// 	if err == nil {
	// 		data := buff[:n]
	// 		if _, err := bin.Write(data); err != nil {
	// 			panic(err)
	// 		}
	// 		continue
	// 	}
	// 	if err == io.EOF {
	// 		if n != 0 {
	// 			data := buff[:n]
	// 			if _, err := bin.Write(data); err != nil {
	// 				panic(err)
	// 			}
	// 		}
	// 		break
	// 	}
	// }

	// fmt.Printf("read %d bytes\n", bin.Len())
}
