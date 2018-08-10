// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package sound

import "errors"

// ChannelAlignmentError is returned by sources and sinks
// when the input/output slices aren't sized correctly
// with respect to the number of channels
// in the sink/source.  More concretely this error is returned
// when:
//
//  Sink.Send(d) is called and len(d) % Sink.Channels() != 0
//
//  Source.Receive(d) is called and len(d) % Source.Channels() != 0
//
// Because the buffers d don't contain a number of samples representing
// one sample for all channels for all represented points in time.
var ChannelAlignmentError = errors.New("Channels misaligned")
