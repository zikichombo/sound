// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package sample

func ToCmplx(dst []complex128, src []float64) []complex128 {
	if cap(dst) < len(src) {
		dst = make([]complex128, len(src))
	}
	dst = dst[:len(src)]
	for i := range dst {
		dst[i] = complex(src[i], 0)
	}
	return dst
}

func FromCmplx(dst []float64, src []complex128) []float64 {
	if cap(dst) < len(src) {
		dst = make([]float64, len(src))
	}
	dst = dst[:len(src)]
	for i := range dst {
		dst[i] = real(src[i])
	}
	return dst
}
