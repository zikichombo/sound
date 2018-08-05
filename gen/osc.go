// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package gen

import "math"

type osc struct {
	a      float64
	y0, y1 float64
}

func newOsc(rps float64) *osc {
	return newOscPh(rps, 0.0)
}

func newOscPh(rps, ph float64) *osc {
	osc := &osc{a: 2.0 * math.Cos(rps)}
	osc.y1 = math.Cos(ph - rps)
	osc.y0 = math.Cos(ph - 2.0*rps)
	return osc
}

func (o *osc) Next() float64 {
	z := o.a*o.y1 - o.y0
	o.y0, o.y1 = o.y1, z
	return z
}
