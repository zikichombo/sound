// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package ops

import (
	"io"

	"github.com/irifrance/snd"
)

type up struct {
	snd.Source
	i, n int
}

func (u *up) Receive(dst []float64) (int, error) {
	nC := u.Channels()
	if len(dst)%nC != 0 {
		return 0, snd.ChannelAlignmentError
	}
	nF := len(dst) / nC
	rnF := (nF - u.i) / u.n
	if (nF-u.i)%u.n != 0 {
		rnF++
	}
	var m int
	var err error
	if rnF > 0 {
		m, err = u.Source.Receive(dst[:rnF*nC])
	}
	if err != nil {
		// handle trailing zeros
		i := 0
		for u.i != 0 {
			for c := 0; c < nC; c++ {
				dst[i] = 0.0
				i++
			}
			u.i++
			if u.i == u.n {
				u.i = 0
			}
		}
		if i == 0 {
			return 0, io.EOF
		}
		return i / nC, nil
	}
	// calc wF, number of output frames
	wF := m * u.n
	if u.i > 0 {
		wF += u.n - u.i
	}
	if wF > nF {
		wF = nF
	}

	for c := nC - 1; c >= 0; c-- {
		sStart := c * m
		dStart := c * wF
		s, t := sStart+m-1, dStart+wF-1
		for t >= dStart {
			if (u.i+(t-dStart))%u.n == 0 {
				dst[t] = dst[s]
				s--
			} else {
				dst[t] = 0.0
			}
			t--
		}
	}
	u.i += wF
	u.i = u.i % u.n
	return wF, nil
}

// UpSample returns an up-sampled version of src.
//
// Specifically, if any channel of src is
// in the form
//
//  s0 s1 s2 ...
//
// Then Up(src, n) is a source whose corresponding
// channels are in the form
//
//  s0 0 ... 0 s1 0 ... 0 s2 0 ... 0 ...
//     -------    -------    -------
//     n times    n times    n times
func Upsample(src snd.Source, n int) snd.Source {
	return &up{Source: src, n: n}
}
