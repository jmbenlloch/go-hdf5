// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hdf5

// #include "blosc_filter.h"
// #include "hdf5.h"
import "C"
import (
	"fmt"
	"unsafe"
)

// Registers Blosc plugin
func RegisterBlosc() (string, string, error) {
	var version *C.char
	var date *C.char
	var r C.int

	/* Register the filter with the library */
	r = C.register_blosc(&version, &date)
	versionGo := C.GoString(version)
	dateGo := C.GoString(date)
	C.free(unsafe.Pointer(version))
	C.free(unsafe.Pointer(date))
	if r != 1 {
		return "", "", fmt.Errorf("Blosc filter failed to register %d", r)
	}
	return versionGo, dateGo, nil
}

type BloscFilter int

const (
	BLOSC_BLOSCLZ BloscFilter = C.BLOSC_BLOSCLZ
	BLOSC_LZ4     BloscFilter = C.BLOSC_LZ4
	BLOSC_LZ4HC   BloscFilter = C.BLOSC_LZ4HC
	BLOSC_SNAPPY  BloscFilter = C.BLOSC_SNAPPY
	BLOSC_ZLIB    BloscFilter = C.BLOSC_ZLIB
	BLOSC_ZSTD    BloscFilter = C.BLOSC_ZSTD
)

type BloscShuffle int

const (
	BLOSC_NOSHUFFLE  BloscShuffle = C.BLOSC_NOSHUFFLE
	BLOSC_SHUFFLE    BloscShuffle = C.BLOSC_SHUFFLE
	BLOSC_BITSHUFFLE BloscShuffle = C.BLOSC_BITSHUFFLE
)

func ConfigureBloscFilter(p_list *PropList, compressionAlgorithm BloscFilter, compressionLevel int, bitShuffle BloscShuffle) error {
	/* This is the easiest way to call Blosc with default values: 5
	   for BloscLZ and shuffle active. */
	/* r = H5Pset_filter(plist, FILTER_BLOSC, H5Z_FLAG_OPTIONAL, 0, NULL); */

	cd_values := make([]uint32, 7)
	/* But you can also taylor Blosc parameters to your needs */
	/* 0 to 3 (inclusive) param slots are reserved. */
	cd_values[4] = uint32(compressionLevel) /* compression level */
	/* 0: shuffle not active, 1: byte-wise shuffle active, 2: bit-wise shuffle */
	cd_values[5] = uint32(bitShuffle)
	//cd_values[6] = C.BLOSC_BLOSCLZ          /* the actual compressor to use */
	//cd_values[6] = C.BLOSC_LZ4 /* the actual compressor to use */
	cd_values[6] = uint32(compressionAlgorithm) /* the actual compressor to use */

	/* Set the filter with 7 params */
	return h5err(C.H5Pset_filter(C.hid_t(p_list.id), C.FILTER_BLOSC, C.H5Z_FLAG_OPTIONAL, 7, (*C.uint)(&cd_values[0])))
}
