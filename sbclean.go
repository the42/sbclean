// Copyright 2011 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

//target:github.com/the42/encoding/sbclean

// sbclean is a package which provides functionality to
//  * encode arbitrary code points into eight bit clean code points
//  * decode previously encoded chunks into originating format
package sbclean

import (
	"io"
	"os"
	"strconv"
)

/*
 * Encoder
 */

// Encode encodes src by removing every most significant bit from every byte
// and appending the accumulating information every eight byte  
// Encode encodes EncodedLen(len(src)) bytes to dst.
func Encode(dst, src []byte) {

	var accu byte
	var sindex, dindex int

	for sindex < len(src) {

		dst[dindex] = src[sindex] >> 1
		accu |= src[sindex] & 1

		dindex++
		sindex++

		if sindex%7 == 0 {
			dst[dindex] = accu
			dindex++
			accu = 0
		} else {
			accu <<= 1
		}
	}

	if sindex%7 != 0 {
		dst[dindex] = accu
	}
}

type encoder struct {
	err  os.Error
	w    io.Writer
	buf  [7]byte    // buffered data waiting to be encoded
	nbuf int        // number of bytes in buf
	out  [1024]byte // output buffer
}

func (e *encoder) Write(p []byte) (n int, err os.Error) {
	if e.err != nil {
		return 0, e.err
	}

	// Leading fringe
	if e.nbuf > 0 {
		var i int
		for i = 0; i < len(p) && e.nbuf < 7; i++ {
			e.buf[e.nbuf] = p[i]
			e.nbuf++
		}
		n += i
		p = p[i:]
		if e.nbuf < 7 {
			return
		}
		Encode(e.out[0:], e.buf[0:])
		if _, e.err = e.w.Write(e.out[0:8]); e.err != nil {
			return n, e.err
		}
		e.nbuf = 0
	}

	// Large interior chunks
	for len(p) >= 7 {
		nn := len(e.out) / 8 * 7
		if nn > len(p) {
			nn = len(p)
		}
		nn -= nn % 7
		if nn > 0 {
			Encode(e.out[0:], p[0:nn])
			if _, e.err = e.w.Write(e.out[0 : nn/7*8]); e.err != nil {
				return n, e.err
			}
		}
		n += nn
		p = p[nn:]
	}

	// Trailing fringe
	for i := 0; i < len(p); i++ {
		e.buf[i] = p[i]
	}
	e.nbuf = len(p)
	n += len(p)
	return
}

// Close flushes any pending output from the encoder.
// It is an error to call Write after calling Close.
func (e *encoder) Close() os.Error {
	// If there's anything left in the buffer, flush it out
	if e.err == nil && e.nbuf > 0 {
		Encode(e.out[0:], e.buf[0:e.nbuf])
		_, e.err = e.w.Write(e.out[0 : e.nbuf+1])
		e.nbuf = 0
	}
	return e.err
}

// Eight bit clean encoding operates in 8-byte blocks; when finished
// writing, the caller must Close the returned encoder to flush any
// partially written blocks.
func NewEncoder(w io.Writer) io.WriteCloser {
	return &encoder{w: w}
}

// EncodedLen returns the length in bytes of the eight bit clean encoding
// of an input buffer of length n.
func EncodedLen(n int) int {
	if n%7 == 0 {
		return n * 8 / 7
	}
	return n*8/7 + 1
}

/*
 * Decoder
 */

type CorruptInputError int64

func (e CorruptInputError) String() string {
	return "expected 8 bit clean data byte but most significant bit is set at " + strconv.Itoa64(int64(e))
}

func Decode(dst, src []byte) (n int, err os.Error) {

	var accu byte
	var index int

	for len(src) > 7 {

		accu = src[7]
		if (accu & 0x80) > 0 {
			return 0, CorruptInputError(n)
		}
		for index = 6; index >= 0; index-- {
			if (src[index] & 0x80) > 0 {
				return 0, CorruptInputError(n)
			}

			dst[index] = src[index] << 1
			dst[index] |= accu & 1
			accu >>= 1
			n++
		}
		src = src[8:]
		dst = dst[7:]
	}

	if len(src) > 0 {
		accu = src[len(src)-1]
		if (accu & 0x80) > 0 {
			return 0, CorruptInputError(n)
		}
		for index = len(src) - 2; index >= 0; index-- {
			if (src[index] & 0x80) > 0 {
				return 0, CorruptInputError(n)
			}
			dst[index] = src[index] << 1
			accu >>= 1
			dst[index] |= accu & 1
			n++
		}
	}

	return n, nil
}

type decoder struct {
	err    os.Error
	r      io.Reader
	buf    [1024]byte // leftover input
	nbuf   int
	out    []byte // leftover decoded output
	outbuf [1024 / 8 * 7]byte
}

func (d *decoder) Read(p []byte) (n int, err os.Error) {
	if d.err != nil {
		return 0, d.err
	}

	// Use leftover decoded output from last read.
	if len(d.out) > 0 {
		n = copy(p, d.out)
		d.out = d.out[n:]
		return n, nil
	}

	// Read a chunk.
	nn := len(p) / 7 * 8
	if nn > len(d.buf) {
		nn = len(d.buf)
	}
	
	// Eight bit clean encoding has no padding; We will read past the end, in which case ErrUnexpectedEOF is not an error
	nn, d.err = io.ReadFull(d.r, d.buf[d.nbuf:nn])
	if d.err != nil && d.err != io.ErrUnexpectedEOF {
		return 0, d.err
	}
	d.nbuf += nn

	// Decode chunk into p, or d.out and then p if p is too small.
	nr := d.nbuf
	nw := d.nbuf / 8 * 7
	if nw > len(p) {
		nw, d.err = Decode(d.outbuf[0:], d.buf[0:nr])
		d.out = d.outbuf[0:nw]
		n = copy(p, d.out)
		d.out = d.out[n:]
	} else {
		n, d.err = Decode(p, d.buf[0:nr])
	}
	d.nbuf -= nr
	for i := 0; i < d.nbuf; i++ {
		d.buf[i] = d.buf[i+nr]
	}

	if d.err == nil {
		d.err = err
	}
	return n, d.err
}

// NewDecoder constructs a new eight bit clean stream decoder.
func NewDecoder(r io.Reader) io.Reader {
	return &decoder{r: r}
}

// DecodedLen returns the maximum length in bytes of the decoded data
// corresponding to n bytes of eight bit clean encoded data.
func DecodedLen(n int) int {
	return n * 7 / 8
}
