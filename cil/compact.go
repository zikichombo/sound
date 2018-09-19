// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package cil

import "zikichombo.org/sound"

// Compact takes a channel de-interleaved slice d with nC channels
// and compacts d in place so that it has nF frames placed channel
// de-interleaved in d[:nF*nC].
//
// If nF >= len(d)/nC, then Compact will panic.
// If len(d)%nC != 0, then Compact will return sound.ErrChannelAlignment.
func Compact(d []float64, nC, nF int) error {
	if nC == 1 {
		return nil
	}
	if len(d)%nC != 0 {
		return sound.ErrChannelAlignment
	}
	dF := len(d) / nC
	if dF == nF {
		return nil
	}
	for c := 0; c < nC; c++ {
		sStart := c * dF
		sEnd := sStart + dF
		dStart := c * nF
		dEnd := dStart + nF
		copy(d[dStart:dEnd], d[sStart:sEnd])
	}
	return nil
}
