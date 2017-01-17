/*
Copyright (C) 2016-2017 Andreas T Jonsson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package xmi

/*
#include "stdlib.h"
#include "xmi2mid.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

const (
	NoConversion    = C.XMIDI_CONVERT_NOCONVERSION
	MT32ToGM        = C.XMIDI_CONVERT_MT32_TO_GM
	MT32ToGS        = C.XMIDI_CONVERT_MT32_TO_GS
	MT32ToGS127     = C.XMIDI_CONVERT_MT32_TO_GS127     // This one is broken, don't use.
	MT32ToGS127Drum = C.XMIDI_CONVERT_MT32_TO_GS127DRUM // This one is broken, don't use.
	GS127ToGS       = C.XMIDI_CONVERT_GS127_TO_GS
)

func ToMidi(xmiFile []byte, convFlag uint32) ([]byte, error) {
	var (
		outputData *C.uint8_t
		outputSize C.uint32_t
	)

	if res := C._WM_xmi2midi((*C.uint8_t)(&xmiFile[0]), C.uint32_t(len(xmiFile)), &outputData, &outputSize, C.uint32_t(convFlag)); res != 0 {
		return nil, fmt.Errorf("conversion error: %v", res)
	}

	ptr := outputData
	data := make([]byte, outputSize)

	for i := 0; i < int(outputSize); i++ {
		data[i] = byte(*ptr)
		ptr = (*C.uint8_t)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + 1))
	}

	C.free(unsafe.Pointer(outputData))
	return data, nil
}
