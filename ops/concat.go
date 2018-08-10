// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package ops

import (
	"io"

	"zikichombo.org/sound"
	"zikichombo.org/sound/freq"
)

type concat struct {
	srcs []sound.Source
	n    int
}

func (c *concat) SampleRate() freq.T {
	return c.srcs[0].SampleRate()
}

func (c *concat) Channels() int {
	return c.srcs[0].Channels()
}

func (c *concat) Receive(dst []float64) (int, error) {
	var err error
	var n, ttl int
	for c.n < len(c.srcs) {
		src := c.srcs[c.n]
		n, err = src.Receive(dst[ttl:])
		ttl += n
		if err == io.EOF {
			c.n++
		}
		if err != nil {
			return 0, err
		}
		if ttl == len(dst) {
			return ttl, nil
		}
	}
	if ttl == 0 {
		return 0, io.EOF
	}
	return ttl, nil
}

func (c *concat) Close() error {
	var err, retErr error
	for i := c.n; i < len(c.srcs); i++ {
		err = c.srcs[i].Close()
		if err != nil && retErr != nil {
			retErr = err
		}
	}
	return retErr
}

// Concat returns the concatenation of ts if ts is non-empty.
// ts should not contain any nil sources, or Concat or the
// returned source may panic.  If ts is empty, Concat returns
// nil.
func Concat(ts ...sound.Source) sound.Source {
	if len(ts) == 0 {
		return nil
	}
	return &concat{
		srcs: ts,
		n:    0}
}
