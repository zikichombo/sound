// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package sound

import "errors"

// Duplex is an interface for full duplex sound I/O.
type Duplex interface {
	// NB: channels returns InChannels() + OutChannels()
	Form
	Closer
	InChannels() int
	OutChannels() int

	// SendReceive on a duplex connection will play out
	// and capture to in.
	//
	// SendReceive returns a ChannelAlignmentError if
	// len(in) is not a multiple of InChannels() and also
	// if len(out) is not a multiple of OutChannels().
	//
	// SendReceive returns a FrameAlignmentError if
	// the number of frames in out is not equal to the number
	// of frames in in.
	//
	// SendReceive returns the number of frames of input, n.
	// n == 0 iff error != nil
	// if n < len(in)/InChannels() then only the first
	// n frames of out are sent and subsequent calls to
	// SendReceive will return 0, io.EOF.
	//
	SendReceive(out, in []float64) (int, error)
}

// See Duplex.SendReceive
var FrameAlignmentError = errors.New("duplex frames misaligned.")
