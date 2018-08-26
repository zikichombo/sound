// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package ops

import (
	"math"

	"zikichombo.org/sound"
	"zikichombo.org/sound/freq"
)

type Envelope interface {
	Amp() float64
}

type expEnv struct {
	count int
	v     float64
	a     float64
}

func (e *expEnv) Amp() float64 {
	e.count++
	e.v *= e.a
	if e.v < 0 {
		e.v = 0
	}
	return e.v
}

type linEnv struct {
	v float64
	r float64
}

func (e *linEnv) Amp() float64 {
	res := e.v
	e.v += e.r
	if e.v < 0 {
		e.v = 0
	}
	if e.v > 1.0 {
		e.v = 1.0
	}
	return res
}

type sinEnv struct {
	v   float64
	rps float64
}

func (e *sinEnv) Amp() float64 {
	res := math.Sin(e.v)
	e.v += e.rps
	return res
}

// ExpEnv returns a source with an exponential envelope
// fading from start with decay decay.
func ExpEnv(src sound.Source, start, decay float64) sound.Source {
	return Env(src, &expEnv{a: decay, v: start})
}

// LinearEnv returns a source with a linear envelope
// fading from start at rate rate.
func LinearEnv(src sound.Source, start, rate float64) sound.Source {
	return Env(src, &linEnv{v: start, r: rate})
}

// SinEnv returns a source with a sin enveope at frequence f.
func SinEnv(src sound.Source, f freq.T) sound.Source {
	return Env(src, &sinEnv{v: 0, rps: src.SampleRate().RadsPer(f)})
}

type eSrc struct {
	sound.Source
	env Envelope
}

func (e *eSrc) Receive(dst []float64) (int, error) {
	nC := e.Source.Channels()
	if len(dst)%nC != 0 {
		return 0, sound.ErrChannelAlignment
	}
	n, err := e.Source.Receive(dst)
	if err != nil {
		return 0, err
	}
	env := e.env
	if nC == 1 {
		for i := 0; i < n; i++ {
			dst[i] *= env.Amp()
		}
		return n, nil
	}
	f := 0
	for f < n {
		a := env.Amp()
		for c := 0; c < nC; c++ {
			dst[c*n+f] *= a
		}
		f++
	}
	return n, nil
}

func Env(src sound.Source, env Envelope) sound.Source {
	return &eSrc{Source: src, env: env}
}
