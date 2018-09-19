// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package ops

import (
	"zikichombo.org/sound"
	"zikichombo.org/sound/freq"
)

type asif struct {
	sound.Source
	f freq.T
}

// AsIf reinterprets a source s as if it were sampled at frequency f.
func AsIf(s sound.Source, f freq.T) sound.Source {
	return &asif{Source: s, f: f}
}

// Freq returns the frequency from the constructor AsIf().
func (a *asif) Freq() freq.T {
	return a.f
}
