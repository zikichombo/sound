// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package gen_test

import (
	"testing"
	"time"

	"zikichombo.org/codec/wav"
	"zikichombo.org/sound/freq"
	"zikichombo.org/sound/gen"
	"zikichombo.org/sound/ops"
)

func TestNote(t *testing.T) {
	a := ops.LimitDur(gen.Note(220*freq.Hertz), time.Second)
	if e := wav.Save(a, "a3.wav"); e != nil {
		t.Fatal(e)
	}

}

func TestNotes(t *testing.T) {
	a3 := 220 * freq.Hertz
	a := ops.LimitDur(gen.Notes(a3, (a3*3)/2, (a3*5)/4), time.Second)
	if e := wav.Save(a, "a3Mj.wav"); e != nil {
		t.Fatal(e)
	}
}
