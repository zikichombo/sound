// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package cil

import "testing"

func TestCompact(t *testing.T) {
	N := 10
	M := 7
	C := 3
	d := make([]float64, N*C)
	for c := 0; c < C; c++ {
		for i := 0; i < M; i++ {
			d[c*N+i] = float64(c*N + i)
		}
	}
	if err := Compact(d, C, M); err != nil {
		t.Fatal(err)
	}
	for c := 0; c < C; c++ {
		for i := 0; i < M; i++ {
			if d[c*M+i] != float64(c*N+i) {
				t.Errorf("got %f not %f\n", d[c*M+i], float64(c*N+i))
			}
		}
	}
}
