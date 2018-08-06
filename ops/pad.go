// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

// Copyright 2017 The IriFrance Audio Authors. All rights reserved.  Use of
// this source code is governed by a license that can be found in the License
// file.

package ops

import (
	"io"

	"zikichombo.org/sound"
)

type pad struct {
	sound.Source
	d   float64
	n   int
	err error
}

// Pad returns a sound.Source ps which contains src
// samples followed by n frames in which each
// sample in each channel has sample p.
func Pad(src sound.Source, p float64, n int) sound.Source {
	return &pad{
		Source: src,
		d:      p,
		n:      n * src.Channels()}
}

func (p *pad) Sample() (float64, error) {
	if p.err != nil {
		return 0, p.err
	}
	if p.n <= 0 {
		p.err = io.EOF
		return 0, io.EOF
	}
	d, e := Sample(p.Source)
	if e == io.EOF {
		p.n--
		return p.d, nil
	}
	if e != nil {
		p.err = e
		return 0, e
	}
	return d, nil
}
