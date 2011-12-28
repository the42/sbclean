sbclean - eight bit clean encoding
===================================

sbclean is a package which provides functionality to

* encode arbitrary code points into eight bit clean code points
* decode previously encoded chunks into originating format

The encoding is fixed length and stream orientated, resulting in an increase in
size of 8 / 7.

Eight bit clean encoding is necessary for legacy systems which assign special
meaning for the most significant bit. Until websockets provide a mean to send
binary data, eight bit clean encoding can be used to encode binary data into an
utf8 string.

Features
--------

The package provides the following functionality:

* Encode / Decode byte strings into eight bit clean byte sequences.
* Read / Write interface for stream operations, cf. [Go encoding
packages](http://golang.org/pkg/encoding/)

Installation
------------

  go install github.com/the42/sbclean

License
-------

The package is released under the [Simplified BSD
License](http://www.freebsd.org/copyright/freebsd-license.html) See file
"LICENSE"


Implementation details
----------------------

Encoding: The most significant bit of a byte meant for encoding is moved at the
least significant position of an accumulating byte. The accumulating byte is
shifted towards the left. Subsequently the most significant bit of the input
byte is cleared and copied to the encoding buffer.

This encoding scheme is repeated six more times or until the input stream is
drained. Afterwards the accumulator carrying the most significant bits of the
cleared input bytes is appended to the encoded stream.

If the output sequence is larger but seven bytes: Every eight byte in the
encoded output sequence is the accumulator.
The last byte of the output sequence is always the accumulating byte.

Decoding: Reversal

Testing
-------

To run the tests:

  go test github.com/the42/sbclean
