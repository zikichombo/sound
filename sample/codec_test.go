// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

// Copyright 2017 The IriFrance Audio Authors. All rights reserved.  Use of
// this source code is governed by a license that can be found in the License
// file.

package sample

import (
	"math/rand"
	"testing"
)

func TestConv(tst *testing.T) {
	for _, c := range Codecs {
		for i := 0; i < 64; i++ {
			m := rand.Float64()
			s := c.FromFloat64(m)
			f := c.ToFloat64(s)
			u := c.FromFloat64(f)
			if s != u {
				tst.Errorf("%s %d -> %f -> %d", c, s, f, u)
			}
		}
	}
}

func TestEncDec(tst *testing.T) {
	for _, c := range Codecs {
		buf := make([]byte, c.Bytes())
		fb := make([]float64, 1)
		for i := 0; i < 1024; i++ {
			m := rand.Float64()
			fb[0] = m
			c.Encode(buf, fb)
			c.Decode(fb, buf)
			v := fb[0]
			c.Encode(buf, fb)
			c.Decode(fb, buf)
			if v != fb[0] {
				tst.Errorf("%s: %f -> %f -> %f\n", c, m, v, fb[0])
			}
		}
	}
}
