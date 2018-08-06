// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

// Copyright 2017 The IriFrance Audio Authors. All rights reserved.  Use of
// this source code is governed by a license that can be found in the License
// file.

package ops

import (
	"io"
	"time"

	"zikichombo.org/sound"
)

type limit struct {
	sound.Source
	n int
}

// Limit returns a source which contains at most n frames
// of samples from s.
func Limit(s sound.Source, n int) sound.Source {
	return &limit{
		Source: s,
		n:      n}
}

// LimitDur limits the source s to at most d duration.
// The exact number of samples allowed is
//
//   floor(d / s.SampleRate().Period())
//
func LimitDur(s sound.Source, d time.Duration) sound.Source {
	f := s.SampleRate()
	p := f.Period()
	n := int(d / p)
	return &limit{
		Source: s,
		n:      n}
}

func (lim *limit) Receive(dst []float64) (int, error) {
	if len(dst) == 0 {
		return 0, nil
	}
	if lim.n <= 0 {
		return 0, io.EOF
	}
	nC := lim.Channels()
	m := len(dst) / nC
	if m > lim.n {
		m = lim.n
	}
	n, err := lim.Source.Receive(dst[:m*nC])
	lim.n -= n
	return n, err
}
