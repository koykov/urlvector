package urlvector

import (
	"reflect"
	"unsafe"
)

func init() {
	// Take raw address of keys source and store it.
	h := (*reflect.SliceHeader)(unsafe.Pointer(&bKeys))
	keysAddr = uint64(h.Data)
}
