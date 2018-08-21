// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package ops

import (
	"fmt"

	"zikichombo.org/sound"
	"zikichombo.org/sound/freq"
)

type join struct {
	srcs []sound.Source
	c    int
}

func (j *join) SampleRate() freq.T {
	return j.srcs[0].SampleRate()
}

func (j *join) Channels() int {
	return len(j.srcs)
}

func (j *join) Receive(dst []float64) (int, error) {
	if len(dst)%len(j.srcs) != 0 {
		return 0, sound.ChannelAlignmentError
	}
	n := len(dst) / len(j.srcs)
	nf := -1
	for i, src := range j.srcs {
		m, err := src.Receive(dst[i*n : (i+1)*n])
		if err != nil {
			return 0, err
		}
		if nf == -1 {
			nf = m
		}
		if m > nf {
			nf = m
		}
		for j := m; j < n; j++ {
			dst[i*n+j] = 0.0
		}
	}
	if nf != n {
		for i := range j.srcs {
			copy(dst[i*nf:(i+1)*nf], dst[i*n:i*n+nf])
		}
	}
	return nf, nil
}

func (j *join) Close() error {
	var err, retErr error
	for _, src := range j.srcs {
		err = src.Close()
		if err != nil {
			retErr = err
		}
	}
	return retErr
}

// Join takes a slice of 1-channel sources and
// joins them into a single len(src)-channel source.
//
// If any of the sources has more than 1 channel,
// or if there is any difference in the sound codec
// or frequency of any of the sources, Join fails
// returning a nil source and a non-nil error.
func Join(srcs ...sound.Source) (sound.Source, error) {
	if len(srcs) == 0 {
		return nil, fmt.Errorf("zero channel source not allowed.")
	}
	if len(srcs) == 1 {
		return srcs[0], nil
	}
	if e := Compat(srcs...); e != nil {
		return nil, e
	}
	return &join{srcs: srcs}, nil
}

// MustJoin is like Join but it panics in case of an
// error.
func MustJoin(srcs ...sound.Source) sound.Source {
	res, e := Join(srcs...)
	if e != nil {
		panic(e.Error())
	}
	return res
}
