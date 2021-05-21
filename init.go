package urlvector

import (
	"reflect"
	"unsafe"
)

func init() {
	h := (*reflect.SliceHeader)(unsafe.Pointer(&bKeys))
	keysAddr = uint64(h.Data)
}
