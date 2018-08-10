// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package sound

// Source is an interface for a source of samples.
type Source interface {
	Form
	Closer
	// Receive places samples in d.
	//
	// Receive returns the number of frames placed in d together with
	// any error.  Receive may use all of d as scatch space.
	//
	// Receive returns a non-nil error if and only if it returns 0 frames
	// received.  Receive may return 0 < n < len(d)/Channels() frames only
	// if the subsequent call will return (0, io.EOF).  As a result, the
	// caller need not be concerned with whether or not the data is "ready".
	//
	// Receive returns multi-channel data in de-interleaved format.
	// If len(d) is not a multiple of Channels(), then Receive returns
	// ChannelAlignmentError.  If Receive returns fewer than len(d)/Channels()
	// frames, then the deinterleaved data of n frames is arranged in
	// the prefix d[:n*Channels()].
	Receive(d []float64) (int, error)
}
