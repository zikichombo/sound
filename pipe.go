// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package sound

import (
	"io"
	"sync"
)

// Pipe creates a pair of Source and Sink such that writes to the sink
// are passed through to reads from the Source.
//
// The returned source, sink are safe for use in multiple goroutines.
func Pipe(v Form) (Source, Sink) {
	pC := make(chan *packet)
	doneC := make(chan struct{})
	nC := make(chan int)
	p := &pipe{Form: v, pC: pC, doneC: doneC, nC: nC}
	return p, p
}

type pipe struct {
	Form
	wMu   sync.Mutex
	once  sync.Once
	pC    chan *packet
	nC    chan int
	doneC chan struct{}
	sl    []float64
}

type packet struct {
	sl []float64
	n  int
}

func (p *pipe) Send(d []float64) error {
	nC := p.Channels()
	if len(d)%nC != 0 {
		return ChannelAlignmentError
	}
	trgFrms := len(d) / nC
	pkt := &packet{sl: d}
	select {
	case <-p.doneC:
		return io.EOF
	default:
		p.wMu.Lock()
		defer p.wMu.Unlock()
	}

	for pkt.n < trgFrms {
		select {
		case p.pC <- pkt:
			n := <-p.nC
			pkt.n += n
		case <-p.doneC:
			return io.EOF
		}
	}
	return nil
}

func (p *pipe) Close() error {
	p.once.Do(func() { close(p.doneC) })
	return nil
}

func (p *pipe) Receive(dst []float64) (int, error) {
	nC := p.Channels()
	if len(dst)%nC != 0 {
		return 0, ChannelAlignmentError
	}
	dPkt := &packet{sl: dst, n: 0}
	inFrms := len(dst) / nC
	select {
	case <-p.doneC:
		return 0, io.EOF
	default:
	}
	for dPkt.n < inFrms {
		select {
		case <-p.doneC:
			if dPkt.n == 0 {
				return 0, io.EOF
			}
			compact(dPkt, nC, inFrms)
			return dPkt.n, nil
		case sPkt := <-p.pC:
			n := copyFrames(dPkt, sPkt, nC)
			dPkt.n += n
			p.nC <- n
		}
	}
	return dPkt.n, nil
}

func compact(pkt *packet, nC, nFrms int) {
	if pkt.n == nFrms {
		return
	}
	n := pkt.n
	sl := pkt.sl
	for c := 0; c < nC; c++ {
		sStart := c * nFrms
		sEnd := sStart + n
		dStart := c * n
		dEnd := dStart + n
		copy(sl[dStart:dEnd], sl[sStart:sEnd])
	}
}

func copyFrames(dst *packet, src *packet, nc int) int {
	if nc == 1 {
		return copy(dst.sl, src.sl)
	}
	sFrms := len(src.sl) / nc
	dFrms := len(dst.sl) / nc
	st := sFrms - src.n
	dt := dFrms - dst.n
	t := st
	if dt < t {
		t = dt
	}
	for c := 0; c < nc; c++ {
		sStart := c*sFrms + src.n
		sEnd := sStart + t
		dStart := c*dFrms + dst.n
		dEnd := dStart + t
		copy(dst.sl[dStart:dEnd], src.sl[sStart:sEnd])
	}
	return t
}
