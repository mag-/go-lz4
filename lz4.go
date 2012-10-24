/*
   lz4.go - Go wrapper for LZ4

   Copyright (C) 2012, Utkan Güngördü <utkan@freeconsole.org>

   BSD 2-Clause License (http://www.opensource.org/licenses/bsd-license.php)

   Redistribution and use in source and binary forms, with or without
   modification, are permitted provided that the following conditions are
   met:

       * Redistributions of source code must retain the above copyright
   notice, this list of conditions and the following disclaimer.
       * Redistributions in binary form must reproduce the above
   copyright notice, this list of conditions and the following disclaimer
   in the documentation and/or other materials provided with the
   distribution.

   THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
   "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
   LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
   A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
   OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
   SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
   LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
   DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
   THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
   (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
   OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

/*
   LZ4 - Fast LZ compression algorithm
   Header File
   Copyright (C) 2011-2012, Yann Collet.
   BSD 2-Clause License (http://www.opensource.org/licenses/bsd-license.php)

   Redistribution and use in source and binary forms, with or without
   modification, are permitted provided that the following conditions are
   met:

       * Redistributions of source code must retain the above copyright
   notice, this list of conditions and the following disclaimer.
       * Redistributions in binary form must reproduce the above
   copyright notice, this list of conditions and the following disclaimer
   in the documentation and/or other materials provided with the
   distribution.

   THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
   "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
   LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
   A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
   OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
   SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
   LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
   DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
   THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
   (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
   OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

   You can contact the author at :
   - LZ4 homepage : http://fastcompression.blogspot.com/p/lz4.html
   - LZ4 source repository : http://code.google.com/p/lz4/
*/

// A Go wrapper for the original lz4 library.
package lz4


/*
#include "lz4.h"
#include "lz4hc.h"

static int _LZ4_compressBound(int isize)  {
	return LZ4_compressBound(isize);
}
*/
import "C"
import "unsafe"

func cast(slice []byte) *C.char {
	return (*C.char)(unsafe.Pointer(&slice[0]))
}

/*
Provides the maximum size that LZ4 may output in a "worst case" scenario (input data not compressible)
primarily useful for memory allocation of output buffer.

isize  : is the input size. Max supported value is ~1.9GB
return : maximum output size in a "worst case" scenario
note : this function is limited by "int" range (2^31-1)
*/
func CompressBound(isize int) int {
	return int(C._LZ4_compressBound(C.int(isize)))
}

/* 
Compresses 'isize' bytes from 'source' into 'dest'.
Destination buffer must be already allocated,
and must be sized to handle worst cases situations (input data not compressible)
Worst case size evaluation is provided by CompressBound()

isize  : is the input size. Max supported value is ~1.9GB
return : the number of bytes written in buffer dest
*/
func Compress(source, dest []byte, isize int) int {
	return int( C.LZ4_compress(cast(source), cast(dest), C.int(isize)) )
}

/*
return : the number of bytes in compressed buffer dest
note : destination buffer must be already allocated. 

To avoid any problem, size it to handle worst cases situations (input data not compressible)
Worst case size evaluation is provided by function CompressBound()
*/
func CompressHighCompression(source, dest []byte, isize int) int {
	return int( C.LZ4_compressHC(cast(source), cast(dest), C.int(isize)) )
}

/*
Compress 'isize' bytes from 'source' into an output buffer 'dest' of maximum size 'maxOutputSize'.
If it cannot achieve it, compression will stop, and result of the function will be zero.
This function never writes outside of provided output buffer.

isize  : is the input size. Max supported value is ~1.9GB
maxOutputSize : is the size of the destination buffer (which must be already allocated)
return : the number of bytes written in buffer 'dest'
or 0 if the compression fails
*/
func CompressLimitedOutput(source, dest []byte, isize, maxOutputSize int) int {
	return int( C.LZ4_compress_limitedOutput(cast(source), cast(dest), C.int(isize), C.int(maxOutputSize)) )
}


/*
osize  : is the output size, therefore the original size
return : the number of bytes read in the source buffer

If the source stream is malformed, the function will stop decoding and return a negative result, indicating the byte position of the faulty instruction
This function never writes beyond dest + osize, and is therefore protected against malicious data packets

note : destination buffer must be already allocated.
its size must be a minimum of 'osize' bytes.
*/
func Decompress(source, dest []byte, osize int) int {
	return int( C.LZ4_uncompress(cast(source), cast(dest), C.int(osize)) )
}

/*
isize  : is the input size, therefore the compressed size
maxOutputSize : is the size of the destination buffer (which must be already allocated)
return : the number of bytes decoded in the destination buffer (necessarily <= maxOutputSize)

If the source stream is malformed, the function will stop decoding and return a negative result, indicating the byte position of the faulty instruction
This function never writes beyond dest + maxOutputSize, and is therefore protected against malicious data packets

note   : Destination buffer must be already allocated.
This version is slightly slower than Decompress()
*/
func DecompressUnknownOutputSize(source, dest []byte, isize, maxOutputSize int) int {
	return int( C.LZ4_uncompress_unknownOutputSize(cast(source), cast(dest), C.int(isize), C.int(maxOutputSize)) )
}