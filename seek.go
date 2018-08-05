// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

// Copyright 2017 The IriFrance Audio Authors. All rights reserved.  Use of
// this source code is governed by a license that can be found in the License
// file.

package sound

import "time"

// Interface Seeker provides common access to a
// source of sound with seek support.
type Seeker interface {
	Form
	Pos() int64         // Where we currently are in the source
	Len() int64         // Length of the source in samples
	Seek(f int64) error // Seek goes to the frame index f
}

// SourceSeeker is a Source with Seek support.
type SourceSeeker interface {
	Seeker
	Closer
	// Receive is as specified in Source.Receive.
	Receive(dst []float64) (int, error)
}

// SinkSeeker is a Sink with Seek support.
type SinkSeeker interface {
	Seeker
	Closer
	// Send is as specified as in Sink.Send.
	Send(src []float64) error
}

// RandomAccess is an interface for full read/write/seek
// support.
type RandomAccess interface {
	Seeker
	// Receive is specified as in Source.
	Receive(dst []float64) (int, error)
	// Send is specified as in Sink.
	Send(src []float64) error
}

func d2s(v Form, d time.Duration) int64 {
	return int64(d / v.SampleRate().Period())
}

func s2d(v Form, f int64) time.Duration {
	return v.SampleRate().Period() * time.Duration(f)
}

// When takes a Seeker and returns its position in terms of time.
func When(s Seeker) time.Duration {
	return s2d(s, s.Pos())
}

// Durations takes a Seeker and returns its total duration in
// time.
func Duration(s Seeker) time.Duration {
	return s2d(s, s.Len())
}

// SeekDur takes a Seeker and seeks to the greatest frame
// not exceeding d time.
//
// SeekDur returns the error given by the call to s.Seek()
// implementing the seek.
func SeekDur(s Seeker, d time.Duration) error {
	return s.Seek(d2s(s, d))
}
