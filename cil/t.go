// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package cil

import "github.com/irifrance/snd"

// Type T encapsulates state for a buffer of
// a fixed number of frames and channels.
type T struct {
	c, f int
	bits []uint64
}

// New creates a new T for flat buffers
// of nfrms frames each containing 1 data point
// for each of nch channels.
func New(nch, nfrm int) *T {
	n := nch * nfrm
	sz := n / 64
	if sz*64 < n {
		sz++
	}
	return &T{
		c:    nch,
		f:    nfrm,
		bits: make([]uint64, sz)}
}

// Frames returns the number of frames associated with bufferes processed by t.
func (t *T) Frames() int {
	return t.f
}

// SetFrames sets the frames size of the [de]interleaver t
// to f.
func (t *T) SetFrames(f int) {
	if f <= t.f {
		t.f = f
		return
	}
	n := t.c * f
	sz := n / 64
	if sz*64 < n {
		sz++
	}
	t.f = f
	if cap(t.bits) < sz {
		t.bits = make([]uint64, sz, (sz*5)/3)
	}
	t.bits = t.bits[:sz]
}

// Channels returns the number of channels associated with bufferes processed by t.
func (t *T) Channels() int {
	return t.c
}

// Len returns the length of a raw buffer.
func (t *T) Len() int {
	return t.c * t.f
}

// Inter takes as input de-interleaved raw data and lays it out so it
// is interleaved, one frame after another, each frame containing 1 sample
// for each channel.
//
// If len(raw) is not a multiple of the number of channels expected by t
// (specified in the constructor), then Inter returns snd.ChannelAlignmentError.
//
// raw's frame count may differ from that given in the constructor.
//
func (t *T) Inter(raw []float64) error {
	if t.c == 1 {
		return nil
	}
	N := len(raw)
	if N%t.c != 0 {
		return snd.ChannelAlignmentError
	}
	orgF := t.f
	defer t.SetFrames(orgF)
	f := N / t.c
	t.SetFrames(f)
	t.clear()
	for i := 0; i < N; i++ {
		if !t.setDone(i) {
			continue
		}
		v := raw[i]
		var w float64
		h := i
		j := t.ToDeinter(i)
		for j != i {
			w = raw[j]
			raw[h] = w
			t.set(j)
			h, j = j, t.ToDeinter(j)
		}
		raw[h] = v
	}
	return nil
}

// Deinter rearranges raw so that it is organized as
// a sequence of per-channel time sequences
//
// raw's frame count may differ from that given in t's constructor.
//
// Deinter returns a snd.ChannelAlignmentError if len(raw)%nC != 0
// where nC is the number of channels given in the constructor.
func (t *T) Deinter(raw []float64) error {
	if t.c == 1 {
		return nil
	}
	N := len(raw)
	if N%t.c != 0 {
		return snd.ChannelAlignmentError
	}
	orgF := t.f
	defer t.SetFrames(orgF)
	f := N / t.c
	t.SetFrames(f)
	t.clear()
	for i := 0; i < N; i++ {
		if !t.setDone(i) {
			continue
		}
		v := raw[i]
		var w float64
		h := i
		j := t.ToInter(i)
		for j != i {
			w = raw[j]
			raw[h] = w
			t.set(j)
			h, j = j, t.ToInter(j)
		}
		raw[h] = v
	}
	return nil
}

// Frame returns the subslice of raw which
// corresponds to frame f.  Frame assumes
// raw is interleaved.
func (t *T) Frame(f int, raw []float64) []float64 {
	return raw[t.c*f : t.c*(f+1)]
}

// Chan returns the subslice of raw which corresponds
// to the sequence of data in channel c.  Chan assumes
// raw is de-interleaved.
func (t *T) Chan(c int, raw []float64) []float64 {
	return raw[t.f*c : t.f*(c+1)]
}

// AppendFrames appends each subslice of raw corresponding
// to a frame to dst.  AppendFrames assumes raw is interleaved.
func (t *T) AppendFrames(dst [][]float64, raw []float64) [][]float64 {
	for i := 0; i < t.f; i++ {
		dst = append(dst, t.Frame(i, raw))
	}
	return dst
}

// AppendChans appends each subslice of raw corresponding
// to a channel to dst.  AppendChans assumes raw is de-interleaved.
func (t *T) AppendChans(dst [][]float64, raw []float64) [][]float64 {
	for i := 0; i < t.c; i++ {
		dst = append(dst, t.Chan(i, raw))
	}
	return dst
}

// InterIdx gives the channel and frame number of index i
// for interleaved data.
func (t *T) InterIdx(i int) (c, f int) {
	return i % t.c, i / t.c
}

// DeinterIdx gives the channel and frame number of index i
// for de-interleaved data.
func (t *T) DeinterIdx(i int) (c, f int) {
	return i / t.f, i % t.f
}

// ToDeinter returns the de-interleaved index which corresponds
// to the interleaved index i.
func (t *T) ToDeinter(i int) int {
	c, f := t.InterIdx(i)
	return c*t.f + f
}

// ToInter returns the interleaved index which corresponds to
// the de-interleaved index i.
func (t *T) ToInter(i int) int {
	c, f := t.DeinterIdx(i)
	return f*t.c + c
}

// InterFrom returns the interleaved index of
// sample f in channel c.
func (t *T) InterFrom(c, f int) int {
	return f*t.c + c
}

// DeinterFrom returns the de-interleaved index
// of sample f in channel c.
func (t *T) DeinterFrom(c, f int) int {
	return c*t.f + f
}

func (t *T) setDone(i int) bool {
	j, m := i/64, uint(i%64)
	v := t.bits[j]
	has := v&(1<<m) != 0
	t.bits[j] |= (1 << m)
	return !has
}

func (t *T) set(i int) {
	j, m := i/64, uint(i%64)
	t.bits[j] |= (1 << m)
}

func (t *T) clear() {
	for i := range t.bits {
		t.bits[i] = 0
	}
}
