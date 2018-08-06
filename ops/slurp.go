// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package ops

import (
	"zikichombo.org/sound"
	"zikichombo.org/sound/sample"
)

// SlurpCmplx is like Receive, but puts the result in
// a complex128 slice.
func SlurpCmplx(src sound.Source, dst []complex128) (int, error) {
	sd := make([]float64, len(dst))
	n, e := src.Receive(sd)
	for i := range dst[:n] {
		dst[i] = complex(sd[i], 0)
	}
	return n, e
}

func SlurpFixed(src sound.Source, dst []int64, bps int) (int, error) {
	tmp := make([]float64, len(dst))
	n, e := src.Receive(tmp)
	N := src.Channels() * n
	sample.ToFixeds(dst[:N], tmp[:N], bps)
	return n, e
}
