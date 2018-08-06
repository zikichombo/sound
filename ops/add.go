// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

// Copyright 2017 The IriFrance Audio Authors. All rights reserved.  Use of
// this source code is governed by a license that can be found in the License
// file.

package ops

import (
	"zikichombo.org/sound"
	"zikichombo.org/sound/freq"
)

// Add mixes the srcs into a single source.
// If any of the sources has different sampling rate
// or codec or number of channels, then Add
// returns a nil source and a non-nil error.
func Add(srcs ...sound.Source) (sound.Source, error) {
	if len(srcs) == 0 {
		return nil, nil
	}
	if e := Compat(srcs...); e != nil {
		return nil, e
	}
	res := &sAdd{srcs: srcs}
	return res, nil
}

// MustAdd is like Add but panics if the srcs are incompatible or 0-length.
func MustAdd(srcs ...sound.Source) sound.Source {
	src, err := Add(srcs...)
	if err != nil {
		panic(err.Error())
	}
	return src
}

type sAdd struct {
	srcs []sound.Source
	buf  []float64
	err  error
}

func (a *sAdd) Receive(dst []float64) (int, error) {
	N := len(dst)
	if cap(a.buf) < N {
		a.buf = make([]float64, (N*5)/3)
	}
	a.buf = a.buf[:N]
	for i := range dst {
		dst[i] = 0.0
	}
	nC := a.Channels()
	frms := -1
	var err error
	var n int
	for _, src := range a.srcs {
		n, err = src.Receive(a.buf)
		if err != nil {
			return 0, err
		}
		if frms == -1 {
			frms = n
		}
		if frms != n {
			return 0, sound.ChannelAlignmentError
		}
		for i := 0; i < nC*frms; i++ {
			dst[i] += a.buf[i]
		}
	}
	return frms, nil
}

func (a *sAdd) SampleRate() freq.T {
	return a.srcs[0].SampleRate()
}

func (a *sAdd) Channels() int {
	return a.srcs[0].Channels()
}

func (a *sAdd) Close() error {
	var err, retErr error
	for _, src := range a.srcs {
		err = src.Close()
		if err != nil && retErr != nil {
			retErr = err
		}
	}
	return retErr
}

type sMul struct {
	sAdd
}

func (m *sMul) Receive(dst []float64) (int, error) {
	N := len(dst)
	if cap(m.buf) < N {
		m.buf = make([]float64, (N*5)/3)
	}
	m.buf = m.buf[:N]
	for i := range dst {
		dst[i] = 0.0
	}
	nC := m.Channels()
	frms := -1
	var err error
	var n int
	for _, src := range m.srcs {
		n, err = src.Receive(m.buf)
		if err != nil {
			return 0, err
		}
		if frms == -1 {
			frms = n
		}
		if frms != n {
			return 0, sound.ChannelAlignmentError
		}
		for i := 0; i < nC*frms; i++ {
			dst[i] *= m.buf[i]
		}
	}
	return frms, nil
}
