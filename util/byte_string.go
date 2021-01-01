// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/12

package util

import (
	"reflect"
	"unsafe"
)

// String cost-free conversion of byte stream to character
func String(b []byte) (s string) {
	pBytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pString := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pString.Data = pBytes.Data
	pString.Len = pBytes.Len
	return
}

// Slice cost-free conversion of characters to byte streams
func Slice(s string) (b []byte) {
	pBytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pString := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pBytes.Data = pString.Data
	pBytes.Len = pString.Len
	pBytes.Cap = pString.Len
	return
}
