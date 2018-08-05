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

func TestDecimateMonoChan(t *testing.T) {
	N := 10
	d := make([]float64, N)
	e := make([]float64, N/5)
	for i := range d {
		d[i] = 2.0*rand.Float64() - 1.0
	}
	src := buf.FromSlice(d, freq.Hertz)
	dec := Decimate(src, 3)
	t.Logf("dec(3): %v:\n", d)
	for {
		n, err := dec.Receive(e)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("\t%v\n", e[:n])
	}
}

func TestDecimateStereo(t *testing.T) {
	N := 14
	d := make([]float64, N)
	e := make([]float64, 4)
	for i := range d {
		d[i] = 2.0*rand.Float64() - 1.0
	}
	src := buf.FromSliceChans(d, 2, freq.Hertz)
	dec := Decimate(src, 3)
	t.Logf("dec(3): %v:\n", d)
	for {
		n, err := dec.Receive(e)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("\t%v\n", e[:n*2])
	}
}
