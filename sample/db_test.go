// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package sample

import (
	"math"
	"math/rand"
	"testing"
)

func TestDb(t *testing.T) {
	for i := 0; i < 64; i++ {
		v := rand.Float64()
		w := FromDb(ToDb(v))
		if math.Abs(v-w) > 0.00000000001 {
			t.Errorf("%f -> %f -> %f\n", v, ToDb(v), w)
		}
	}
}
