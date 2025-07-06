package util

import (
	"fmt"
	"strconv"
)

type Intish interface {
	int | int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64
}

// Itoh is helpful for dumping "hex" representations of ints
func Itoh[I Intish](i I) string {
	return fmt.Sprintf("%x", i)
}

// ItohN is helpful for dumping "hex" representations of N count ints
func ItohN[I Intish](iN []I) string {
	s := ""
	for _, i := range iN {
		s += Itoh(i) + " "
	}

	return s
}

// Btoh is helpful for dumping "hex" representations of bytes (uint8)
// if dealing with a sized byte array, you can create a slice that uses the underlying array (no copy)
// by performing Btoh(b[:])
func Btoh(b []byte) string {
	s := ""
	for _, bb := range b {
		s += fmt.Sprintf("%x", bb)
	}
	return s
}

// HStoi is hex-string to i
func HStoi(s string) (int64, error) {
	return strconv.ParseInt(s, 16, 16)
}
