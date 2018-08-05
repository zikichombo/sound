// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

// Copyright 2017 The IriFrance Audio Authors. All rights reserved.  Use of
// this source code is governed by a license that can be found in the License
// file.

package gen

import (
	"github.com/irifrance/snd"
	"github.com/irifrance/snd/freq"
	"github.com/irifrance/snd/ops"
)

var _t = Default()

func Sin(f freq.T) snd.Source {
	return _t.Sin(f)
}

func SinPh(f freq.T, off float64) snd.Source {
	return _t.SinPh(f, off)
}

func Sins(fs ...freq.T) snd.Source {
	ss := make([]snd.Source, len(fs))
	for i := range ss {
		ss[i] = Sin(fs[i])
	}
	res, _ := ops.Add(ss...)
	return res
}

func Impulse() snd.Source {
	done := false
	return mkGen(
		func() (float64, bool) {
			if !done {
				done = true
				return 1.0, false
			}
			return 0.0, false
		})
}

func Constant(v float64) snd.Source {
	return mkGen(func() (float64, bool) {
		return v, false
	})
}

func Squares(f freq.T) snd.Source {
	return _t.Squares(f)
}

func Chirp(l, step freq.T) snd.Source {
	return _t.Chirp(l, step)
}

func Noise() snd.Source {
	return _t.Noise()
}

func Silence() snd.Source {
	return mkGen(func() (float64, bool) {
		return 0.0, false
	})
}

func Spikes(fs freq.T) snd.Source {
	return _t.Spikes(fs)
}

func Slice(d []float64) snd.Source {
	return _t.Slice(d)
}

func SliceCmplx(d []complex128) snd.Source {
	return _t.SliceCmplx(d)
}

func Note(f freq.T) snd.Source {
	return _t.Note(f)
}

func Notes(fs ...freq.T) snd.Source {
	return _t.Notes(fs...)
}

func mkGen(f func() (float64, bool)) snd.Source {
	return &s{T: *_t, fn: f}
}
