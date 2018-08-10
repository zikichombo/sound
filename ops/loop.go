// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package ops

import (
	"io"

	"zikichombo.org/sound"
	"zikichombo.org/sound/cil"
)

type loop struct {
	sound.SourceSeeker
	dil *cil.T
	n   int
}

func (l *loop) Receive(dst []float64) (int, error) {
	nC := l.Channels()
	if len(dst)%nC != 0 {
		return 0, sound.ChannelAlignmentError
	}
	if l.n == 0 {
		return 0, io.EOF
	}
	nF := len(dst) / nC
	n, err := l.SourceSeeker.Receive(dst)
	if err != nil && err != io.EOF {
		return 0, err
	}
	if err == io.EOF {
		l.n--
		if l.n == 0 {
			if n == 0 {
				return 0, io.EOF
			}
			return n, nil
		}
		if err := l.Seek(0); err != nil {
			return 0, err
		}
	}
	if n == nF {
		// fast path, "normal" case one buffer
		return nF, nil
	}
	// interleave so append works easily
	// (nb Inter/Deinter are basically free for mono channel anyway)
	l.dil.Inter(dst[:n*nC])
	f := n
	for f < nF {
		n, err = l.SourceSeeker.Receive(dst[f*nC:])
		if err != nil && err != io.EOF {
			return 0, err
		}
		if err == io.EOF {
			l.n--
			if l.n == 0 {
				break
			}
			if err = l.Seek(0); err != nil {
				return 0, err
			}
		}
		l.dil.Inter(dst[f*nC : (f+n)*nC])
		f += n
	}
	l.dil.Deinter(dst[:f*nC])
	return f, nil
}

// Loop creates a src which is src looped n times.
// Loop will loop infinitely if n < 0 and will
// terminate with EOF immediately if n == 0.
func Loop(src sound.SourceSeeker, n int) sound.Source {
	return &loop{SourceSeeker: src, n: n, dil: cil.New(src.Channels(), 1024)}
}
