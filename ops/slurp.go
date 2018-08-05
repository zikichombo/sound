// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

// Copyright 2017 The IriFrance Audio Authors. All rights reserved.  Use of
// this source code is governed by a license that can be found in the License
// file.

package ops

import (
	"github.com/irifrance/snd"
	"github.com/irifrance/snd/sample"
)

// SlurpCmplx is like Receive, but puts the result in
// a complex128 slice.
func SlurpCmplx(src snd.Source, dst []complex128) (int, error) {
	sd := make([]float64, len(dst))
	n, e := src.Receive(sd)
	for i := range dst[:n] {
		dst[i] = complex(sd[i], 0)
	}
	return n, e
}

func SlurpFixed(src snd.Source, dst []int64, bps int) (int, error) {
	tmp := make([]float64, len(dst))
	n, e := src.Receive(tmp)
	N := src.Channels() * n
	sample.ToFixeds(dst[:N], tmp[:N], bps)
	return n, e
}
