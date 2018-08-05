// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package ops

import "github.com/irifrance/snd"

type decimate struct {
	snd.Source
	i, n int
	buf  []float64
}

func (d *decimate) rFrames(f int) int {
	i := d.i
	j, k := 0, 0
	for j < f {
		i++
		if i == d.n {
			i = 0
			k++
			continue
		}
		j++
		k++
	}
	return k
}

func (d *decimate) wFrames(f int) int {
	i := d.i
	j, k := 0, 0
	for k < f {
		i++
		if i == d.n {
			i = 0
			k++
			continue
		}
		j++
		k++
	}
	return j
}

func (d *decimate) Receive(dst []float64) (int, error) {
	nC := d.Channels()
	if len(dst)%nC != 0 {
		return 0, snd.ChannelAlignmentError
	}
	nF := len(dst) / nC
	rF := d.rFrames(nF)
	N := rF * nC
	if cap(d.buf) < N {
		d.buf = make([]float64, N, (N*5)/3)
	}
	d.buf = d.buf[:N]
	n, err := d.Source.Receive(d.buf)
	if err != nil {
		return 0, err
	}
	wF := d.wFrames(n)
	i, j := 0, 0
	for i < wF {
		d.i++
		if d.i == d.n {
			d.i = 0
			j++
			continue
		}
		for c := 0; c < nC; c++ {
			dst[c*wF+i] = d.buf[c*n+j]
		}
		i++
		j++
	}
	return wF, nil
}

// Decimate returns an n-decimated src, i.e. Decimate drops
// every n'th frame from src and puts the resulting stream
// in the form of a snd.Source.
func Decimate(src snd.Source, n int) snd.Source {
	return &decimate{Source: src, n: n}
}
