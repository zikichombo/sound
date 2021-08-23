// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package ops

import "github.com/zikichombo/sound"

type zeros struct {
	sound.Source
	l, r  float64
	z     float64
	count int
}

func (z *zeros) Sample() (float64, error) {
	d, e := Sample(z.Source)
	if e != nil {
		return d, e
	}
	if d == 0.0 {
		return d, nil
	}
	z.l, z.r = z.r, d
	l, r := z.l, z.r
	if l < r && l < 0.0 && 0.0 < r {
		z.count++
	} else if r < l && 0.0 < l && r < 0.0 {
		z.count++
	}
	return d, nil
}

// Zeros counts the zero crossings in src.  src should be
// mono-channel.
func Zeros(src sound.Source) int {
	zs := &zeros{Source: src}
	for {
		_, e := zs.Sample()
		if e != nil {
			return zs.count
		}
	}
}
