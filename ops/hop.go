// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package ops

import "github.com/irifrance/snd"

// Hop shifts buf left by shift elements and reads
// up to shift elements into buf at the end by means
// of src.Receive.
func Hop(src snd.Source, buf []float64, shift int) (int, error) {
	copy(buf, buf[shift:])
	return src.Receive(buf[len(buf)-shift:])
}
