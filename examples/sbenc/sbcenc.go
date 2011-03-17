// Copyright 2011 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

// Seven bit clean encode or decode data. Operates on stdin and stdout.
package main

import (
	"fmt"
	"flag"
	"github.com/the42/encoding/sbclean"
	"os"
	"io"
)


func main() {

	var encode, decode, helpsi bool

	flag.BoolVar(&encode, "e", false, "seven bit encode stdin to stdout. default")
	flag.BoolVar(&decode, "d", false, "seven bit decode stdin to stdout")
	flag.BoolVar(&helpsi, "h", false, "print this help screen")

	flag.Parse()

	if helpsi || encode && decode {
		fmt.Println("\nSeven bit encoder\n")
		fmt.Println("Usage: sbencode [-e|-d] for encode or decode of seven bit clean data")
		fmt.Println("\t-t prints this help")
		os.Exit(1)
	}

	if !(encode || decode) {
		encode = true
	}

	var err os.Error

	if encode {
		encoder := sbclean.NewEncoder(os.Stdout)
		_, err = io.Copy(encoder, os.Stdin)
		encoder.Close()
	} else {
		decoder := sbclean.NewDecoder(os.Stdin)
		_, err = io.Copy(os.Stdout, decoder)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
		os.Exit(2)
	}
}
