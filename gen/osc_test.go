// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package gen

import (
	"math"
	"testing"

	"zikichombo.org/sound/freq"
)

const eps = 1e-10

func TestOsc(t *testing.T) {
	sf := 48000 * freq.Hertz
	of := 100 * freq.Hertz
	rps := sf.RadsPer(of)
	o := newOsc(rps)
	for i := 0; i < 1<<14; i++ {
		ov, cv := o.Next(), math.Cos(float64(i)*rps)
		if math.Abs(ov-cv) > eps {
			t.Errorf("%d: osc %0.12f cos %0.12f\n", i, ov, cv)
		}
	}
	/*
		op := newOscPh(rps, math.Pi/2)
		for i := 0; i < 512; i++ {
			ov, cv := op.Next(), math.Cos(float64(i)*rps+math.Pi/2)
			if math.Abs(ov-cv) > eps {
				t.Errorf("osc %0.12f cos %0.12f\n", ov, cv)
			}
		}
	*/
}
