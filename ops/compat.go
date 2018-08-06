// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

// Copyright 2017 The IriFrance Audio Authors. All rights reserved.  Use of
// this source code is governed by a license that can be found in the License
// file.

package ops

import (
	"fmt"

	"zikichombo.org/sound"
)

// Compat checks that srcs are compatible in
// terms of sample codec, number of channels,
// and sample rate.
func Compat(srcs ...sound.Source) error {
	if len(srcs) == 0 {
		return nil
	}
	ref := srcs[0]
	C := ref.Channels()
	f := ref.SampleRate()
	for i := 1; i < len(srcs); i++ {
		s := srcs[i]
		if s.Channels() != C {
			return fmt.Errorf("incompatible sources %d %d channels", C, s.Channels())
		}
		if s.SampleRate() != f {
			return fmt.Errorf("incompatible sources frequency %s, %s", f, s.SampleRate())
		}
	}
	return nil
}
