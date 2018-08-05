// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package ops

import "github.com/irifrance/snd"

type amp struct {
	snd.Source
	by float64
}

func (a *amp) Receive(dst []float64) (int, error) {
	n, e := a.Source.Receive(dst)
	for i := 0; i < n; i++ {
		dst[i] *= a.by
	}
	return n, e
}

func Amplify(src snd.Source, by float64) snd.Source {
	return &amp{Source: src, by: by}
}