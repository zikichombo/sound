// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package sound

import "fmt"
import "zikichombo.org/sound/freq"

type _form struct {
	c int
	f freq.T
}

// Channels returns the number of channels of data
// in form v.
func (v *_form) Channels() int {
	return v.c
}

// SampleRate returns the sample rate of data passing through v.
func (v *_form) SampleRate() freq.T {
	return v.f
}

func (v *_form) String() string {
	return fmt.Sprintf("[%d channels @ %s]", v.c, v.f)
}

// NewForm creates a new Form for sampling frequency f
// with c channels of data.
func NewForm(f freq.T, c int) Form {
	return &_form{c: c, f: f}
}

// Interface Form specifies the logical content of pcm audio data,
// namely the number of channels and the sample rate.  Form does not specify
// the in-memory layout of pcm audio data such as sample codec or whether or
// not the data is channel-interleaved.
//
// Applications which normalize memory layout of pcm data can use Form
// implementations to determine all necessary values of pcm data.
type Form interface {
	Channels() int
	SampleRate() freq.T
}

// MonoCd returns a single channel form at CD sampling rate.
func MonoCd() Form {
	return &_form{c: 1, f: 44100 * freq.Hertz}
}

// MonoDvd gives a mono-channel form at Dvd sample rate (48kHz).
func MonoDvd() Form {
	return &_form{c: 1, f: 48000 * freq.Hertz}
}

// StereoCd returns a stereo form at CD sampling rate.
func StereoCd() Form {
	return &_form{c: 2, f: 44100 * freq.Hertz}
}

// StereoDvd gives a 2-channel form at Dvd sample rate (48kHz).
func StereoDvd() Form {
	return &_form{c: 2, f: 48000 * freq.Hertz}
}
