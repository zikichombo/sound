// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package ops

import (
	"testing"

	"zikichombo.org/sound/freq"
	"zikichombo.org/sound/sndbuf"
)

func TestLoopMono(t *testing.T) {
	d := make([]float64, 10)
	for i := 0; i < 5; i++ {
		d[i] = -1.0
	}
	for i := 5; i < 10; i++ {
		d[i] = 1.0
	}
	src := sndbuf.FromSlice(d, freq.Hertz)
	loop := Loop(src, 3)
	e := make([]float64, 25)
	n, err := loop.Receive(e)
	if err != nil {
		t.Fatal(err)
	}
	if n != 25 {
		t.Errorf("got %d not 25\n", n)
	}
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			v := e[5*i+j]
			if i%2 == 0 {
				if v != -1.0 {
					t.Errorf("got %f not %f\n", v, -1.0)
				}
			} else {
				if v != 1.0 {
					t.Errorf("got %f not %f\n", v, 1.0)
				}
			}
		}
	}
}

func TestLoopMonoSmallRecv(t *testing.T) {
	d := make([]float64, 10)
	for i := 0; i < 5; i++ {
		d[i] = -1.0
	}
	for i := 5; i < 10; i++ {
		d[i] = 1.0
	}
	src := sndbuf.FromSlice(d, freq.Hertz)
	loop := Loop(src, 3)
	e := make([]float64, 1)
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			n, err := loop.Receive(e)
			if err != nil {
				t.Fatal(err)
			}
			if n != 1 {
				t.Errorf("got %d not 1\n", n)
			}
			v := e[0]
			if i%2 == 0 {
				if v != -1.0 {
					t.Errorf("got %f not %f\n", v, -1.0)
				}
			} else {
				if v != 1.0 {
					t.Errorf("got %f not %f\n", v, 1.0)
				}
			}
		}
	}
}

func TestLoopStereo(t *testing.T) {
	d := make([]float64, 20)
	for i := 0; i < 5; i++ {
		d[2*i] = -1.0
		d[2*i+1] = -1.0
	}
	for i := 5; i < 10; i++ {
		d[2*i] = 1.0
		d[2*i+1] = 1.0
	}
	src := sndbuf.FromSliceChans(d, 2, freq.Hertz)
	loop := Loop(src, 3)
	e := make([]float64, 50)
	n, err := loop.Receive(e)
	if err != nil {
		t.Fatal(err)
	}
	if n != 25 {
		t.Errorf("got %d not 25\n", n)
	}
	for off := 0; off < 50; off += 25 {
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				v := e[off+5*i+j]
				if i%2 == 0 {
					if v != -1.0 {
						t.Errorf("got %f not %f\n", v, -1.0)
					}
				} else {
					if v != 1.0 {
						t.Errorf("got %f not %f\n", v, 1.0)
					}
				}
			}
		}
	}
}
