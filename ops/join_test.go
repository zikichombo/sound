// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package ops

import (
	"io"
	"testing"

	"zikichombo.org/sound/freq"
	"zikichombo.org/sound/sndbuf"
)

func TestJoin(t *testing.T) {
	ones := make([]float64, 130)
	for i := range ones {
		ones[i] = 1
	}
	twos := make([]float64, 130)
	for i := range twos {
		twos[i] = 2
	}
	one := sndbuf.FromSlice(ones, 44100*freq.Hertz)
	two := sndbuf.FromSlice(twos, 44100*freq.Hertz)
	src, _ := Join(one, two)

	d := make([]float64, 64)
	ttl := 0
	for {
		n, e := src.Receive(d)
		if e == io.EOF {
			break
		}
		ttl += n
		for i := 0; i < n; i++ {
			if d[i] != 1.0 {
				t.Errorf("at %d got %f not %f\n", i, d[i], 1.0)
			}
		}
		for i := n; i < 2*n; i++ {
			if d[i] != 2.0 {
				t.Errorf("at %d got %f not %f\n", i, d[i], 2.0)
			}
		}
	}
	if ttl != 130 {
		t.Errorf("ttl %d not 130\n", ttl)
	}
}
