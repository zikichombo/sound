// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package ops

import (
	"time"

	"github.com/zikichombo/sound"
)

type after struct {
	sound.Source
	n int
}

func (a *after) Receive(dst []float64) (int, error) {
	var m, n int
	var e error
	for a.n > 0 {
		m = len(dst)
		if m > a.n {
			m = a.n
		}
		n, e = a.Source.Receive(dst[:m])
		if e != nil {
			return 0, e
		}
		a.n -= n
	}
	return a.Source.Receive(dst)
}

func After(src sound.Source, n int) sound.Source {
	return &after{Source: src, n: n}
}

func AfterDur(src sound.Source, d time.Duration) sound.Source {
	p := src.SampleRate().Period()
	n := d / p
	return After(src, int(n))
}
