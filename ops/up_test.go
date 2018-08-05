// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package ops

import (
	"io"
	"math/rand"
	"testing"

	"github.com/irifrance/snd/buf"
	"github.com/irifrance/snd/freq"
)

func TestUpSampleMono(t *testing.T) {
	N := 10
	d := make([]float64, N)
	e := make([]float64, N)
	for i := range d {
		e[i] = 2.0*rand.Float64() - 1.0
	}
	t.Logf("input: %v\n", e)
	src := buf.FromSlice(e, freq.Hertz)
	up := Upsample(src, 3)
	n, err := up.Receive(d)
	if n != N {
		t.Errorf("expected %d got %d\n", N, n)
	}
	t.Logf("output chunk 1: %v\n", d[:n])
	n, err = up.Receive(d)
	if n != N {
		t.Errorf("expected %d got %d\n", N, n)
	}
	t.Logf("output chunk 2: %v\n", d[:n])
	n, err = up.Receive(d)
	if n != N {
		t.Errorf("expected %d got %d\n", N, n)
	}
	t.Logf("output chunk 3: %v\n", d[:n])
	n, err = up.Receive(d)
	if err != io.EOF {
		t.Errorf("expected EOF, got %d\n", n)
	}
}

func TestUpSampleStereo(t *testing.T) {
	N := 8
	d := make([]float64, N)
	e := make([]float64, N)
	for i := range e {
		e[i] = 2.0*rand.Float64() - 1.0
	}
	t.Logf("input: %v\n", e)
	src := buf.FromSliceChans(e, 2, freq.Hertz)
	up := Upsample(src, 3)
	n, err := up.Receive(d)
	if n*2 != N {
		t.Errorf("expected %d got %d\n", N, n)
	}
	t.Logf("output chunk 1: %v\n", d[:n*2])
	n, err = up.Receive(d)
	if n*2 != N {
		t.Errorf("expected %d got %d\n", N, n)
	}
	t.Logf("output chunk 2: %v\n", d[:n*2])
	n, err = up.Receive(d)
	if n*2 != N {
		t.Errorf("expected %d got %d\n", N, n)
	}
	t.Logf("output chunk 3: %v\n", d[:n*2])
	n, err = up.Receive(d)
	if err != io.EOF {
		t.Errorf("expected EOF, got %d\n", n)
	}
}
