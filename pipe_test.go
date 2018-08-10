// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package sound

import (
	"io"
	"testing"
)

func TestPipe(t *testing.T) {
	mono := MonoCd()
	_ = mono
	testPipe(t, mono, 1024, 1024)
	testPipe(t, mono, 256, 1024)
	testPipe(t, mono, 1024, 1)
	stereo := StereoCd()
	testPipe(t, stereo, 1024, 1024)
	testPipe(t, stereo, 3, 4)
	testPipe(t, stereo, 1024, 1)
}

func testPipe(t *testing.T, valve Form, n, m int) {
	nC := valve.Channels()
	r, w := Pipe(valve)
	wb := make([]float64, n*nC)
	rb := make([]float64, m*nC)
	for i := range wb {
		wb[i] = float64(i / n)
	}
	go func() {
		for i := 0; i < 128; i++ {
			if err := w.Send(wb); err != nil {
				t.Fatal(err)
			}
		}
		w.Close()
	}()
	ttl := 0
	for ttl < n*128 {
		n, e := r.Receive(rb)
		if n != len(rb)/nC && e == nil {
			t.Fatal("nothing read, no error")
		}
		if e == io.EOF {
			if n != 0 {
				t.Errorf("read %d with eof\n", n)
			}
			break
		}
		if n != len(rb)/nC {
			t.Errorf("read %d not %d\n", n, len(rb)/valve.Channels())
		}
		rEnd := n * nC
		for i := range rb[:rEnd] {
			if int(rb[i]) != i/n {
				t.Errorf("at %d got %d not %d\n", i, int(rb[i]), i/n)
			}
		}
		ttl += n
	}
	if ttl != n*128 {
		t.Errorf("ttl got %d not %d\n", ttl, n*128)
	}
}
