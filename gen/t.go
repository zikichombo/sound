// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

// Copyright 2017 The IriFrance Audio Authors. All rights reserved.  Use of
// this source code is governed by a license that can be found in the License
// file.

package gen

import (
	"io"
	"math"
	"math/rand"

	"zikichombo.org/sound"
	"zikichombo.org/sound/freq"
	"zikichombo.org/sound/ops"
	"zikichombo.org/sound/sample"
)

type T struct {
	freq freq.T
	fn   func(float64) float64
}

func Default() *T {
	return &T{freq: 44100 * freq.Hertz}
}

func New(f freq.T) *T {
	return &T{freq: f}
}

func (t *T) SampleRate() freq.T {
	return t.freq
}

func (t *T) Channels() int {
	return 1
}

type s struct {
	T
	fn func() (float64, bool)
}

func (s *s) Close() error {
	return nil
}

func (s *s) Receive(dst []float64) (int, error) {
	for i := range dst {
		d, e := s.Sample()
		if e == io.EOF {
			return i, nil
		}
		if e != nil {
			return 0, e
		}
		dst[i] = d
	}
	return len(dst), nil
}

func (s *s) Sample() (float64, error) {
	v, e := s.fn()
	if e {
		return 0, io.EOF
	}
	return v, nil
}

func (t *T) Sin(f freq.T) sound.Source {
	return t.SinPh(f, 0)
}

func (t *T) SinPh(f freq.T, off float64) sound.Source {
	var n float64
	rps := t.freq.RadsPer(f)
	fn := func() (float64, bool) {
		v := math.Sin(rps*n + off)
		n += 1.0
		return v, false
	}
	return &s{
		T:  *t,
		fn: fn}

}

func (t *T) Cos(f freq.T) sound.Source {
	return t.CosPh(f, 0)
}

func (t *T) CosPh(f freq.T, off float64) sound.Source {
	var n float64
	rps := t.freq.RadsPer(f)
	fn := func() (float64, bool) {
		v := math.Cos(rps*n + off)
		n += 1.0
		return v, false
	}
	return &s{
		T:  *t,
		fn: fn}
}

func (t *T) Chirp(l, step freq.T) sound.Source {
	instF := t.freq.RadsPer(l)
	stepRads := t.freq.RadsPer(step)
	n := 0.0
	f := func() (float64, bool) {
		v := math.Sin(instF * n)
		n++
		instF += stepRads
		return v, false
	}
	return &s{T: *t, fn: f}
}

func (t *T) Sins(rpss ...freq.T) sound.Source {
	ss := make([]sound.Source, len(rpss))
	for i := range ss {
		ss[i] = t.Sin(rpss[i])
	}
	res, _ := ops.Add(ss...)
	return res
}

func (t *T) Impulse() sound.Source {
	done := false
	return &s{T: *t,
		fn: func() (float64, bool) {
			if !done {
				done = true
				return 1.0, false
			}
			return 0.0, false
		}}
}

func (t *T) Constant(v float64) sound.Source {
	return &s{T: *t,
		fn: func() (float64, bool) {
			return v, false
		}}
}

func (t *T) Squares(f freq.T) sound.Source {
	rps := t.freq.RadsPer(f)
	cur := float64(0)
	return &s{T: *t, fn: func() (float64, bool) {
		res := 1.0
		if cur >= math.Pi {
			res = -1.0
		}
		cur += rps
		if cur >= 2*math.Pi {
			cur -= 2 * math.Pi
		}
		return res, false
	}}
}

func (t *T) Noise() sound.Source {
	return &s{T: *t, fn: func() (float64, bool) {
		return rand.Float64(), false
	}}
}

func (t *T) Silence() sound.Source {
	return &s{T: *t, fn: func() (float64, bool) {
		return 0.0, false
	}}
}

func (t *T) Spikes(f freq.T) sound.Source {
	rps := t.freq.RadsPer(f)
	cur := rps
	return &s{T: *t, fn: func() (float64, bool) {
		cur += rps
		if cur >= 2*math.Pi {
			cur -= 2 * math.Pi
			return 1, false
		}
		return 0, false
	}}
}

func (t *T) Slice(d []float64) sound.Source {
	i := 0
	return &s{T: *t, fn: func() (float64, bool) {
		if i == len(d) {
			return 0, true
		}
		v := d[i]
		i++
		return v, false
	}}
}

func (t *T) SliceCmplx(d []complex128) sound.Source {
	sl := sample.FromCmplx(nil, d)
	return t.Slice(sl)
}

func (t *T) Note(f freq.T) sound.Source {
	alpha := 0.8
	N := 16
	sins := make([]sound.Source, 0, N)

	for i := 0; i < N; i++ {
		a := math.Pow(alpha, float64(i))
		var s sound.Source
		sf := f * freq.T(i+1)
		if t.freq.RadsPer(sf) >= math.Pi {
			break
		}
		s = ops.Amplify(t.SinPh(f*freq.T(i+1), rand.Float64()*2*math.Pi), a/float64(N))
		sins = append(sins, s)
	}
	r, e := ops.Add(sins...)
	if e != nil {
		panic(e.Error())
	}
	return r
}

func (t *T) Notes(f ...freq.T) sound.Source {
	ns := make([]sound.Source, len(f))
	for i := range ns {
		ns[i] = ops.Amplify(t.Note(f[i]), 1/float64(len(ns)))
	}
	r, _ := ops.Add(ns...)
	return r
}
