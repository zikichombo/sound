// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package ops

import "github.com/zikichombo/sound"

func Sample(src sound.Source) (float64, error) {
	var b [1]float64
	_, e := src.Receive(b[:])
	return b[0], e
}
