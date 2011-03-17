// Copyright 2011 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

// Automated tests for the sbclean package
package sbclean

import (
	"os"
	"testing"
	"io/ioutil"
	"io"
	"bytes"
)

// ## sbencode & decode


var files = []string{"test.dat", "test1.dat", "test2.dat"}

func TestSbencodeDecodefromExtFile(t *testing.T) {

	for _, filename := range files {
		f, err1 := os.Open(filename, os.O_RDONLY, 0666)
		if err1 != nil {
			t.Fatalf("Error opening file %s: %s", filename, err1)
		}
		defer f.Close()

		readdata, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatalf("Error reading file: %s", err)
		}

		buffer1 := bytes.NewBuffer(nil)
		buffer2 := bytes.NewBuffer(nil)

		if err != nil {
			t.Fatalf("Error encoding: %s", err)
		}

		enc := NewEncoder(buffer1)
		_, err = io.Copy(enc, bytes.NewBuffer(readdata))
		enc.Close()

		dec := NewDecoder(buffer1)
		_, err = io.Copy(buffer2, dec)
		if err != nil {
			t.Fatalf("Error decoding: %s", err)
		}

		if bytes.Compare(readdata, buffer2.Bytes()) != 0 {
			t.Fatalf("Encode / decode mismatch for file %s", filename)
		}
	}
}
