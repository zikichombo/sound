// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package freq

import (
	"fmt"
	"math"
)

// Type I represents a frequency interval.
type I float64

const (
	Octave   = I(2)
	Semitone = I(1.059463094359)
	Cent     = I(1.000577789507)
)

// From applies the interval i to f returning the result.
func (i I) From(f T) T {
	r := float64(i) * float64(f)
	if r > 0 {
		return T(int64(math.Floor(r + 0.5)))
	} else {
		return T(int64(math.Floor(r - 0.5)))
	}
}

// NFrom applies the interval i to f n times returning the result.
func (i I) NFrom(f T, n int) T {
	for j := 0; j < n; j++ {
		f = i.From(f)
	}
	return f
}

// N returns the interval n times larger than i.
func (i I) N(n int) I {
	r := i
	for j := 0; j < n; j++ {
		r *= i
	}
	return r
}

// Within returns whether or not the frequencies a,b are withing i
// of each other.
func (i I) Within(a, b T) bool {
	if a > b {
		a, b = b, a
	}
	return i.From(a) >= b
}

// Diff gives the (signed) interval from ref to off.
func Diff(ref, off T) I {
	return I(float64(off) / float64(ref))
}

func (i I) parts() (o, s, c int) {
	for i >= Octave {
		i /= 2
		o++
	}
	for i >= Semitone {
		i /= Semitone
		s++
	}
	for i >= (Cent - 0.0000000001) {
		i /= Cent
		c++
	}
	return
}

func (i I) Octaves() int {
	j := 0
	for i >= Octave {
		i /= 2
		j++
	}
	return j
}

func (i I) Semitones() int {
	_, s, _ := i.parts()
	return s
}

func (i I) Cents() int {
	_, _, c := i.parts()
	return c
}

func (i I) String() string {
	return fmt.Sprintf("[%d o %d s %d c]", i.Octaves(), i.Semitones(), i.Cents())
}

// Temper gives an I for octaves divided eually into n components.
func Temper(n int) I {
	v := math.Pow(2, 1/float64(n))
	return I(v)
}
