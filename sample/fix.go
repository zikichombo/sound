// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package sample

func ToFixed(d float64, nBits int) int64 {
	d = float64(float32(d))
	q := float64(int64(1 << uint(nBits-1)))
	s := 1.0 / q
	r := int64(d / s)
	return r
}

func ToFixeds(dst []int64, src []float64, nBits int) []int64 {
	if cap(dst) < len(src) {
		dst = make([]int64, len(src))
	}
	dst = dst[:len(src)]
	for i, v := range src {
		dst[i] = ToFixed(v, nBits)
	}
	return dst
}

func ToFloat(d int64, nBits int) float64 {
	s := float64(int64(1 << uint(nBits-1)))
	return float64(d) / s
}

func ToFloats(dst []float64, src []int64, nBits int) []float64 {
	if cap(dst) < len(src) {
		dst = make([]float64, len(src))
	}
	dst = dst[:len(src)]
	for i, v := range src {
		dst[i] = ToFloat(v, nBits)
	}
	return dst
}
