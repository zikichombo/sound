// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package ops

import (
	"fmt"
	"io"

	"zikichombo.org/sound"
)

// Copy copies samples from src to dst until
// EOF.
//
// Copy returns non-nil error if and only if:
//
// 1. src returns non EOF error on Receive(); or
//
// 2. dst return non-nil error on Send(); or
//
// 3. src and dst are incompatible.
//
func Copy(dst sound.Sink, src sound.Source) error {
	if src.SampleRate() != dst.SampleRate() || src.Channels() != dst.Channels() {
		return fmt.Errorf("incompatible source/sink for copy: %s v %s\n", sound.Form(src), sound.Form(dst))
	}
	nC := src.Channels()
	buf := make([]float64, 1024*nC)
	var e error
	var n int
	for {
		n, e = src.Receive(buf)
		if e == io.EOF {
			return nil
		}
		if e != nil {
			return e
		}
		if e = dst.Send(buf[:n*nC]); e != nil {
			return e
		}
	}
	return nil
}
