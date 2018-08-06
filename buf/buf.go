// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

// Copyright 2017 The IriFrance Audio Authors. All rights reserved.  Use of
// this source code is governed by a license that can be found in the License
// file.

package buf

import (
	"fmt"
	"io"

	"zikichombo.org/sound"
	"zikichombo.org/sound/cil"
	"zikichombo.org/sound/freq"
)

// Type Buf implements an in memory sound Seeker/Source/Sink
type T struct {
	dat   []float64
	freq  freq.T
	nchan int
	pos   int
	il    *cil.T
}

// FromSource reads in the source into a buffer (*T).
// If an error occurs while reading other than io.EOF,
// that error is returned.
func FromSource(src sound.Source) (*T, error) {
	b := New(src.SampleRate(), src.Channels())
	buf := make([]float64, 1024)
	for {
		n, e := src.Receive(buf)
		if e == io.EOF {
			break
		}
		if e != nil {
			return nil, e
		}
		b.Send(buf[:n*src.Channels()])
	}
	b.Seek(0)
	return b, nil
}

// New creates a new snd buffer at sampling frequency f with c channels.
func New(f freq.T, c int) *T {
	return &T{
		dat:   make([]float64, 0, 1024),
		freq:  f,
		nchan: c,
		pos:   0}
}

// FromSlice creates a single channel buffer backed by ds.
//
// FromSlice assumes ds is mono-channel.
func FromSlice(ds []float64, f freq.T) *T {
	return FromSliceChans(ds, 1, f)
}

// FromSliceChans creates a buffer from slice ds in channel interleaved
// format with samplerate f.
//
// FromSliceChans panics if len(ds) is not a multiple of nc
func FromSliceChans(ds []float64, nc int, f freq.T) *T {
	if len(ds)%nc != 0 {
		panic("channel alignment")
	}
	return &T{
		dat:   ds,
		freq:  f,
		nchan: nc,
		pos:   0}
}

// Slice returns the slice of in-memory samples storing the snd data.
//
// The backing slice is channel-interleaved.  This differs from
// the Source/Sink expected interface; however it is a natural fit
// for appendable data.  To deinterleave the result, see the
// cil package.
func (b *T) Slice() []float64 {
	return b.dat
}

// Split returns a slice of buffers, one per channel.
func (b *T) Split() []*T {
	return b.SplitTo(nil)
}

// SplitTo places one buffer for each channel in
// dst and returns it.  If dst doesn't have sufficient
// capacity, a new slice is returned in its place.
func (b *T) SplitTo(dst []*T) []*T {
	if cap(dst) < b.nchan {
		dst = make([]*T, b.nchan)
	}
	dst = dst[:b.nchan]
	if b.il == nil {
		b.il = cil.New(b.nchan, len(b.dat)/b.nchan)
	}
	b.il.Deinter(b.dat)
	for i := 0; i < b.nchan; i++ {
		csl := make([]float64, len(b.dat)/b.nchan)
		copy(csl, b.il.Chan(i, b.dat))
		dst[i] = FromSlice(csl, b.SampleRate())
	}
	b.il.Inter(b.dat)
	return dst
}

// SampleRate returns the sampling frequency of the snd buffer.
func (b *T) SampleRate() freq.T {
	return b.freq
}

// Channels returns the number of channels in the snd buffer.
func (b *T) Channels() int {
	return b.nchan
}

// Receive implements snd.Source, placing channel-deinterleaved
// data in dst.
func (b *T) Receive(dst []float64) (int, error) {
	if len(dst)%b.nchan != 0 {
		return 0, sound.ChannelAlignmentError
	}

	n := len(dst)
	m := len(b.dat) - b.pos
	if m == 0 {
		return 0, io.EOF
	}
	if m > n {
		m = n
	}
	frms := m / b.nchan
	c := 0
	f := 0
	// nb b.pos must be at a frame boundary, enforced
	// by pared-down interface
	for i := b.pos; i < b.pos+m; i++ {
		dst[c*frms+f] = b.dat[i]
		c++
		if c == b.nchan {
			c = 0
			f++
		}
	}
	b.pos += m
	return f, nil
}

// Send implements sound.Sink, taking channel-deinterleaved data in d
// and placing it in the memory buffer.
func (b *T) Send(d []float64) error {
	if len(d)%b.nchan != 0 {
		return sound.ChannelAlignmentError
	}

	frms := len(d) / b.nchan
	f := 0
	c := 0
	for f < frms {
		si := c*frms + f
		b.dat = append(b.dat, d[si])
		c++
		if c == b.nchan {
			c = 0
			f++
		}
	}
	return nil
}

func (b *T) Close() error {
	return nil
}

// Len returns the length of the buffer in terms of frames.
func (b *T) Len() int64 {
	return int64(len(b.dat) / b.nchan)
}

// Pos returns the current position in the buffer in terms of
// frames.
func (b *T) Pos() int64 {
	return int64(b.pos / b.nchan)
}

// Seek seeks to the frame with index f.  Seek returns a non-nil out
// of range error iff f < 0 or f > b.Len().
//
func (b *T) Seek(f int64) error {
	if f < 0 || f > b.Len() {
		return fmt.Errorf("io error: out of range")
	}
	b.pos = int(f) * b.nchan
	return nil
}
