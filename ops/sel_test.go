// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package ops

import (
	"math/rand"
	"testing"

	"github.com/zikichombo/sound"
	"github.com/zikichombo/sound/freq"
)

func TestSel(t *testing.T) {
	N := 8
	nC := 2
	sC := 1
	d := make([]float64, N*nC)
	for i := range d {
		d[i] = rand.Float64() - 0.5
	}
	src, snk := sound.Pipe(sound.NewForm(44100*freq.Hertz, nC))
	defer snk.Close()
	go snk.Send(d)
	selSrc := Select(src, sC)
	e := make([]float64, N)
	n, err := selSrc.Receive(e)
	if err != nil {
		t.Fatal(err)
	}
	if n != N {
		t.Fatalf("only read %d/%d\n", n, N)
	}
	for i, v := range e {
		j := sC*N + i
		if d[j] != v {
			t.Errorf("output %d in %d got %f not %f\n", i, j, v, d[j])
		}
	}
}
