// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package sample

import (
	"math/rand"
	"testing"
)

func TestFix(t *testing.T) {
	N := 1024
	for i := 0; i < N; i++ {
		d := rand.Float64()
		nb := rand.Intn(16) + 16
		fixed := ToFixed(d, nb)
		flt := ToFloat(fixed, nb)
		fix2 := ToFixed(flt, nb)
		if fix2 != fixed {
			t.Errorf("bit loss %d %d\n", fixed, fix2)
		}
	}
}
