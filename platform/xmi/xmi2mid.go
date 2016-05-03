/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

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
