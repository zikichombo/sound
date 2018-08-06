// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package freq

import (
	"math"
	"testing"
	"time"
)

func TestSPC(t *testing.T) {
	sf := 192 * KiloHertz
	f := 1920 * Hertz
	sps := sf.SamplesPerCycle(f)
	if math.Abs(sps-100) > 0.0000000001 {
		t.Errorf("sps of %s at %s gave %f\n", sf, f, sps)
	}
	g := sf.FromSamplesPerCycle(sps)
	if f != g {
		t.Errorf("%s -> %s\n", f, g)
	}
}

func TestFreqPeriod(t *testing.T) {
	f := 2 * Hertz
	if f.Period() != time.Second/2 {
		t.Errorf("wrong period %s != %s\n", f.Period(), time.Second/2)
	}
	if FromPeriod(time.Second/2) != f {
		t.Errorf("from period %s didn't give %s, gave %s\n", time.Second/2, f, FromPeriod(time.Second/2))
	}
	f += 500 * MilliHertz
	if f.Period() != (2*time.Second)/5 {
		t.Errorf("wrong period %s != %s\n", f.Period(), (2*time.Second)/5)
	}
	if FromPeriod(time.Second) != Hertz {
		t.Errorf("from period %s gave %s\n", time.Second, FromPeriod(time.Second))
	}

}

type fct struct {
	F T
	D time.Duration
	C int
}

func TestFreqCycles(t *testing.T) {
	tsts := []fct{
		{440 * Hertz, time.Second, 440},
		{1 * KiloHertz, time.Second, 1000},
		{MilliHertz, 5000 * time.Second, 5}}
	for _, tst := range tsts {
		r, _ := tst.F.Cycles(tst.D)
		if r != tst.C {
			t.Errorf("wrong number of cycles for (%s, %s): %d != %d", tst.F, tst.D, r, tst.C)
		}
	}
}

func TestFreqPhase(t *testing.T) {
	p := Hertz.Phase(time.Second / 2)
	if math.Abs(math.Pi-p) > 0.0001 {
		t.Errorf("phase wrong: %f vs %f(pi)", p, math.Pi)
	}
}

func TestRadiansPerSample(t *testing.T) {
	r100 := 2.0 * math.Pi / 100
	r1000 := 2.0 * math.Pi / 1000
	rate := 44100 * Hertz
	f := 441 * Hertz
	if f.RadsPerAt(rate) != r100 {
		t.Errorf("%f unexpected, wanted %f", f.RadsPerAt(rate), r100)
	}
	rate -= 100 * Hertz
	f = 44 * Hertz
	if f.RadsPerAt(rate) != r1000 {
		t.Errorf("%f unexpected, wanted %f", f.RadsPerAt(rate), r1000)
	}

	rate += 100 * Hertz
	f = 440 * Hertz
	rA4 := 2.0 * math.Pi / (100.0 + 10.0/44.0)

	if f.RadsPerAt(rate) != rA4 {
		t.Errorf("%s at %s: %f unexpected, wanted %f", f, rate, f.RadsPerAt(rate), rA4)
	}
}

func TestFreqOf(t *testing.T) {
	s := 44100 * Hertz
	f := s / 100
	rps := f.RadsPerAt(s)
	g := s.FreqOf(rps)
	if f != g {
		t.Errorf("differing freqs %s %s", f, g)
	}
}
