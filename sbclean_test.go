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
func TestSbencodeDecodefromExtFile(t *testing.T) {

var cnt int64	
  f, err1 := os.Open("test.dat", os.O_RDONLY, 0666)
	if err1 != nil {
		t.Fatalf("Error opening file test.dat: %s", err1)
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
	cnt, err = io.Copy(enc, bytes.NewBuffer(readdata))
	
	print(buffer1.Len())
	enc.Close()

	dec := NewDecoder(buffer1)
	cnt, err = io.Copy(buffer2, dec)
		if err != nil {
		t.Fatalf("Error decoding: %s", err)
	}
	print(cnt)
	
	i1 := buffer2.Len()
	i2 := len(readdata)
	if i1 != i2 {
	  t.Fatalf("Size mismatch: %d - %d", i1, i2)
	}

}
