// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package ops

import (
	"zikichombo.org/sound"
	"zikichombo.org/sound/freq"
)

type sel struct {
	src sound.Source
	c   int
	buf []float64
}

// Select returns a Source, let us call it sc, corresponding to channel c
// of s.
//
// The caller should use sc in place of s as sc.Receive(d) simply reads from
// s, discarding unused channels and placing data from the selected channel in d.
//
// Select panics if c is out of bounds w.r.t. s.Channels().
func Select(s sound.Source, c int) sound.Source {
	if s.Channels() == 1 {
		return s
	}
	return &sel{src: s, c: c, buf: make([]float64, s.Channels()*1024)}
}

func (s *sel) SampleRate() freq.T {
	return s.src.SampleRate()
}

func (s *sel) Channels() int {
	return 1
}

func (s *sel) Close() error {
	s.buf = nil
	return s.src.Close()
}

func (s *sel) Receive(dst []float64) (int, error) {
	nC := s.src.Channels()
	nF := len(dst)
	var err error
	var n, f int
	for f < nF {
		end := len(s.buf) / nC
		if end > nF-f {
			end = nF - f
		}
		n, err = s.src.Receive(s.buf[:end*nC])
		if err != nil {
			return 0, err
		}
		copy(dst[f:], s.buf[s.c*n:(s.c+1)*n])
		f += n
	}
	return nF, nil
}
