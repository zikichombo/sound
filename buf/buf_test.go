// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

// Copyright 2017 The IriFrance Audio Authors. All rights reserved.  Use of
// this source code is governed by a license that can be found in the License
// file.

package buf_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/irifrance/snd"
	"github.com/irifrance/snd/buf"
	"github.com/irifrance/snd/freq"
)

func TestBufBasic(t *testing.T) {
	b := buf.New(44100*freq.Hertz, 1)
	N := 1024
	for i := 0; i < N; i++ {
		b.Send([]float64{float64(i)})
	}
	if b.Len() != int64(N) {
		t.Errorf("wrong length %d != %d\n", b.Len(), N)
	}
	b.Seek(0)
	buf := make([]float64, 1)
	for i := 0; i < N; i++ {

		_, e := b.Receive(buf)
		if e != nil {
			t.Errorf("unexpected error: %s", e)
		}
		s := buf[0]
		if s != float64(i) {
			t.Errorf("didn't get back %d", i)
		}
	}
	if snd.Duration(b) != time.Duration(b.Len())*b.SampleRate().Period() {
		t.Errorf("unexpected duration %s != %s", snd.Duration(b), b.SampleRate().Period())
	}
	fmt.Printf("dur: %s\n", snd.Duration(b))
	_, ok := interface{}(b).(snd.Source)
	if !ok {
		t.Errorf("doesn't fit interface.")
	}
}
