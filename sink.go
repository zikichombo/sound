// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package sound

import "zikichombo.org/sound/freq"

// Sink is an interface for a destination of samples.
type Sink interface {
	Form
	Closer
	// Send sends all the samples in d, returning a non-nil
	// error if it does not succeed.
	//
	// len(d) must be a multiple of the number of channels
	// associated with the Sink, or Send will return a
	// ChannelAlignmentError.
	//
	// In case the Sink is multi-channel, d is interpreted in
	// channel-deinterleaved format.
	Send(d []float64) error
}

// Discard does nothing when samples are sent to it.
var Discard Sink = &discard{}

type discard struct{}

func (*discard) Send(d []float64) error {
	return nil
}

func (*discard) Close() error {
	return nil
}

func (*discard) SampleRate() freq.T {
	return 0
}

func (*discard) Channels() int {
	return 1
}
