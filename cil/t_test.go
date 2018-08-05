// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package cil_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/irifrance/snd/cil"
)

func TestCilIdx(t *testing.T) {
	C := 8
	F := 1024
	for i := 0; i < 128; i++ {
		c := rand.Intn(C) + 1
		f := rand.Intn(F) + 1
		cilr := cil.New(c, f)
		for i := 0; i < c; i++ {
			for j := 0; j < f; j++ {
				n := j*c + i
				ii, jj := cilr.InterIdx(n)
				if ii != i {
					t.Fatal("inter c")
				}
				if jj != j {
					t.Fatal("inter f")
				}
				n = i*f + j
				ii, jj = cilr.DeinterIdx(n)
				if ii != i {
					t.Fatal("deinter c")
				}
				if jj != j {
					t.Fatal("dinter f")
				}
			}
		}
	}
}

func TestCilInterDeinter(t *testing.T) {
	C := 8
	F := 1024
	for i := 0; i < 128; i++ {
		c := rand.Intn(C) + 1
		f := rand.Intn(F) + 1
		fmt.Printf("testing %d channels and %d frames\n", c, f)
		testCilCF(c, f, t)
	}
}

func testCilCF(c, f int, t *testing.T) {
	n := c * f
	raw := make([]float64, n)
	tmp := make([]float64, n)
	for i := range raw {
		raw[i] = rand.Float64()
	}
	copy(tmp, raw)
	cilr := cil.New(c, f)
	cilr.Deinter(raw)
	cilr.Inter(raw)
	for i := range raw {
		if raw[i] != tmp[i] {
			t.Fatalf("inter(deinter(raw)) at %d", i)
		}
	}
	cilr.Inter(raw)
	cilr.Deinter(raw)
	for i := range raw {
		if raw[i] != tmp[i] {
			t.Fatalf("deinter(inter(raw)) at %d", i)
		}
	}
}
