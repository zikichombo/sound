// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package sample

import (
	"encoding/binary"
	"math"
)

type Codec int

const (
	SInt8 Codec = iota
	SByte
	SInt16L   // little endian 2 byte
	SInt16B   // big endian 2 byte
	SInt24L   // little endian 3 byte
	SInt24B   // big endian 3 byte
	SInt32L   // little endian 4 byte
	SInt32B   // big endian 4 byte
	SFloat32L // float32 bits stored uint32 little endian
	SFloat32B // float32 bits stored uint32 big endian
	SFloat64L // float64 bits stored uint64 little endian
	SFloat64B // float64 bits stored uint64 big endian
)

// Codecs lists all codecs.
var Codecs = [...]Codec{
	SInt8, SByte, SInt16L, SInt16B, SInt24L, SInt24B, SInt32L, SInt32B,
	SFloat32L, SFloat32B, SFloat64L, SFloat64B}

func (s Codec) String() string {
	switch s {
	case SInt8:
		return "SInt8"
	case SByte:
		return "SByte"
	case SInt16L:
		return "SInt16L"
	case SInt16B:
		return "SInt16B"
	case SInt24L:
		return "SInt24L"
	case SInt24B:
		return "SInt24B"
	case SInt32L:
		return "SInt32L"
	case SInt32B:
		return "SInt32B"
	case SFloat32L:
		return "SFloat32L"
	case SFloat32B:
		return "SFloat32B"
	case SFloat64L:
		return "SFloat64L"
	case SFloat64B:
		return "SFloat64B"
	default:
		panic("unnown codec")
	}
}

// IsFloat returns whether s is a floating point codec.
func (s Codec) IsFloat() bool {
	return s >= SFloat32L
}

// ByteOrder returns the byte order of s.
//
// Special cases: SInt8 and SByte return binary.LittleEndian.
func (s Codec) ByteOrder() binary.ByteOrder {
	switch s {
	case SInt8, SByte, SInt16L, SInt32L, SInt24L, SFloat32L, SFloat64L:
		return binary.LittleEndian
	default:
		return binary.BigEndian
	}
}

// Bits returns the number of bits uses to represent each sample
// with codec s.
func (s Codec) Bits() int {
	switch s {
	case SInt8:
		return 8
	case SByte:
		return 8
	case SInt16L, SInt16B:
		return 16
	case SInt24L, SInt24B:
		return 24
	case SInt32L, SInt32B:
		return 32
	case SFloat32L, SFloat32B:
		return 32
	case SFloat64L, SFloat64B:
		return 64
	default:
		panic("unknown sample type")
	}
}

// Bytes returns the number of bytes used to represent samples with
// codec s.
func (s Codec) Bytes() int {
	return s.Bits() / 8
}

// ToFloat gives a float64 representation of samples with codec s
// whose bits are represented in "bits".
func (s Codec) ToFloat64(bits uint64) float64 {
	var buf [8]byte
	var fb [1]float64
	bo := s.ByteOrder()
	bo.PutUint64(buf[:], bits)
	s.Decode(fb[:], buf[:])
	return fb[0]
}

// FromFloat64 gives the bits of the sample v as a uint64
// encoded using codec s.
func (s Codec) FromFloat64(v float64) uint64 {
	var buf [1]float64
	var dst [8]byte
	bo := s.ByteOrder()
	buf[0] = v
	s.Encode(dst[:], buf[:])
	return bo.Uint64(dst[:])
}

// Decode decodes source samples encoded into src using codec
// s into a slice of float64s.
//
// If len(dst)*s.Bytes() > len(src) then Decode panics.
func (s Codec) Decode(dst []float64, src []byte) {
	le, be := binary.LittleEndian, binary.BigEndian
	switch s {
	case SInt8:
		for i := range dst {
			b := src[i]
			dst[i] = float64(b) / 128.0
		}
	case SByte:
		for i := range dst {
			b := src[i]
			dst[i] = float64(-int8(b)) / 128.0
		}
	case SInt16L:
		start := 0
		for i := range dst {
			dst[i] = float64(int16(le.Uint16(src[start:start+2]))) / float64(1<<15)
			start += 2
		}
	case SInt16B:
		start := 0
		for i := range dst {
			dst[i] = float64(int16(be.Uint16(src[start:start+2]))) / float64(1<<15)
			start += 2
		}
	case SInt24L:
		start := 0
		for i := range dst {
			v := int32(src[start]) | int32(src[start+1])<<8 | int32(src[start+2])<<16
			if v&(1<<23) != 0 {
				v &= (1 << 23) - 1
				v = -v
			}
			dst[i] = float64(v) / float64(1<<23)
			start += 3
		}
	case SInt24B:
		start := 0
		for i := range dst {
			v := int32(src[start+2]) | int32(src[start+1])<<8 | int32(src[start])<<16
			if v&(1<<23) != 0 {
				v &= (1 << 23) - 1
				v = -v
			}
			dst[i] = float64(v) / float64(1<<23)
			start += 3
		}

	case SInt32L:
		start := 0
		for i := range dst {
			dst[i] = float64(int32(le.Uint32(src[start:start+4]))) / float64(1<<31)
			start += 4
		}
	case SInt32B:
		start := 0
		for i := range dst {
			dst[i] = float64(int32(be.Uint32(src[start:start+4]))) / float64(1<<31)
			start += 4
		}
	case SFloat32L:
		start := 0
		for i := range dst {
			dst[i] = float64(math.Float32frombits(le.Uint32(src[start : start+4])))
			start += 4
		}
	case SFloat32B:
		start := 0
		for i := range dst {
			dst[i] = float64(math.Float32frombits(be.Uint32(src[start : start+4])))
			start += 4
		}
	case SFloat64L:
		start := 0
		for i := range dst {
			dst[i] = math.Float64frombits(le.Uint64(src[start : start+8]))
			start += 8
		}
	case SFloat64B:
		start := 0
		for i := range dst {
			dst[i] = math.Float64frombits(be.Uint64(src[start : start+8]))
			start += 8
		}
	default:
		panic("unknown Codec")
	}
}

// Encode encodes samples from src into dst using codec s.
//
// If len(dst) < len(src)*s.Bytes(), then Encode panics.
func (s Codec) Encode(dst []byte, src []float64) {
	le, be := binary.LittleEndian, binary.BigEndian
	switch s {
	case SInt8:
		for i, v := range src {
			dst[i] = byte(int8(v * float64(1<<7)))
		}
	case SByte:
		for i, v := range src {
			dst[i] = byte(-int8(v * float64(1<<7)))
		}
	case SInt16L:
		start := 0
		for _, v := range src {
			le.PutUint16(dst[start:start+2], uint16(v*float64(1<<15)))
			start += 2
		}
	case SInt16B:
		start := 0
		for _, v := range src {
			be.PutUint16(dst[start:start+2], uint16(v*float64(1<<15)))
			start += 2
		}
	case SInt24L:
		start := 0
		for _, v := range src {
			s := int32(0)
			if v < 0 {
				v = -v
				s = 1 << 23
			}
			w := int32(float64(1<<23)*v) | s
			dst[start] = byte(w)
			dst[start+1] = byte(w >> 8)
			dst[start+2] = byte(w >> 16)
			start += 3
		}
	case SInt24B:
		start := 0
		for _, v := range src {
			s := int32(0)
			if v < 0 {
				v = -v
				s = 1 << 23
			}
			w := int32(float64(1<<23)*v) | s
			dst[start+2] = byte(w)
			dst[start+1] = byte(w >> 8)
			dst[start] = byte(w >> 16)
			start += 3
		}

	case SInt32L:
		start := 0
		for _, v := range src {
			le.PutUint32(dst[start:start+4], uint32(v*float64(1<<31)))
			start += 4
		}
	case SInt32B:
		start := 0
		for _, v := range src {
			be.PutUint32(dst[start:start+4], uint32(v*float64(1<<31)))
			start += 4
		}
	case SFloat32L:
		start := 0
		for _, v := range src {
			le.PutUint32(dst[start:start+4], math.Float32bits(float32(v)))
			start += 4
		}
	case SFloat32B:
		start := 0
		for _, v := range src {
			be.PutUint32(dst[start:start+4], math.Float32bits(float32(v)))
			start += 4
		}
	case SFloat64L:
		start := 0
		for _, v := range src {
			le.PutUint64(dst[start:start+8], math.Float64bits(v))
			start += 8
		}
	case SFloat64B:
		start := 0
		for _, v := range src {
			be.PutUint64(dst[start:start+8], math.Float64bits(v))
			start += 8
		}
	default:
		panic("unknown Codec")
	}
}
