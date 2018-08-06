// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package freq

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

// Type Frequency represents a frequency.
type T int64

// Different levels of "ouch".
const (
	NanoHertz  T = 1
	MicroHertz   = 1000 * NanoHertz
	MilliHertz   = 1000 * MicroHertz
	Hertz        = 1000 * MilliHertz
	KiloHertz    = 1000 * Hertz
	MegaHertz    = 1000 * KiloHertz
	GigaHertz    = 1000 * MegaHertz
)

// String implements Stringer.
func (f T) String() string {
	if f < 0 {
		if -f >= 0 {
			return fmt.Sprintf("-%s", -f)
		}
		return "-Inf Hz"
	}
	if f < MicroHertz {
		return fmt.Sprintf("%dnHz", f)
	}
	if f < MilliHertz {
		if f%MicroHertz == 0 {
			return fmt.Sprintf("%.3f\u00B5Hz", float64(f)/float64(MicroHertz))
		} else {
			return fmt.Sprintf("%.6f\u00B5Hz", float64(f)/float64(MicroHertz))
		}
	}
	if f < Hertz {
		if f%MilliHertz == 0 {
			return fmt.Sprintf("%.3fmiHz", float64(f)/float64(MilliHertz))
		} else if f%MicroHertz == 0 {
			return fmt.Sprintf("%.6fmiHz", float64(f)/float64(MilliHertz))
		} else {
			return fmt.Sprintf("%.9fmiHz", float64(f)/float64(MilliHertz))
		}
	}
	if f < KiloHertz {
		if f%Hertz == 0 {
			return fmt.Sprintf("%.0fHz", float64(f)/float64(Hertz))
		} else if f%MilliHertz == 0 {
			return fmt.Sprintf("%.3fHz", float64(f)/float64(Hertz))
		} else if f%MicroHertz == 0 {
			return fmt.Sprintf("%.6fHz", float64(f)/float64(Hertz))
		} else {
			return fmt.Sprintf("%.9fHz", float64(f)/float64(Hertz))
		}
	}
	if f < MegaHertz {
		if f%KiloHertz == 0 {
			return fmt.Sprintf("%.0fkHz", float64(f)/float64(KiloHertz))
		} else if f%Hertz == 0 {
			return fmt.Sprintf("%.3fkHz", float64(f)/float64(KiloHertz))
		} else if f%MilliHertz == 0 {
			return fmt.Sprintf("%.6fkHz", float64(f)/float64(KiloHertz))
		} else if f%MicroHertz == 0 {
			return fmt.Sprintf("%.9fkHz", float64(f)/float64(KiloHertz))
		} else {
			return fmt.Sprintf("%.12fkHz", float64(f)/float64(KiloHertz))
		}
	}
	if f < GigaHertz {
		return fmt.Sprintf("%.3fmgHz", float64(f)/float64(MegaHertz))
	}
	return fmt.Sprintf("%.3fgHz", float64(f)/float64(GigaHertz))
}

// FromPeriod gives the frequency whose period is p.
func FromPeriod(p time.Duration) T {
	return T((1000000000 * time.Second) / p)
}

// Float64 give f in Hertz.
func (f T) Float64() float64 {
	return float64(f) / float64(Hertz)
}

// Period gives the period of time of 1 cycle at frequency f.
func (f T) Period() time.Duration {
	return (1000000000 * time.Second) / time.Duration(f)
}

// Cycles gives the number of times a signal at frequency f cycles
// during d time and a remainder duration.
//
// For example if f = 3 Hz and t=.5sec then f.Cycles() returns
// (1, .5sec - .33333 sec).
func (f T) Cycles(d time.Duration) (int, time.Duration) {
	p := f.Period()
	return int(d / p), d % p
}

// Phase gives the offet or phase in radians at time d,
// assuming a signal at frequency f starts at 0.
func (f T) Phase(d time.Duration) float64 {
	p := float64(f.Period())
	df := float64(d)
	return 2.0 * math.Pi * (math.Mod(df, p) / p)
}

// RadPerAt gives the radians per sample of a signal at frequency f sampled at
// frequency r.
func (f T) RadsPerAt(s T) float64 {
	d := float64(f) / float64(s)
	return 2.0 * math.Pi * d
}

// RadsPer gives the radians per sample of g at sampling frequency f.
func (f T) RadsPer(g T) float64 {
	return g.RadsPerAt(f)
}

func (f T) SamplesPerCycle(g T) float64 {
	return float64(f) / float64(g)
}

func (f T) FromSamplesPerCycle(spc float64) T {
	return T(int64(math.Floor(float64(f)/spc + 0.5)))
}

// FreqOf gives the frequency f' such that f' has r radians per sample
// at sample freq f.
func (s T) FreqOf(rps float64) T {
	// find freq f s.t. s.RadPer(f) = rps
	// rps = (2*pi) f/s
	m := rps * float64(s) / (2 * math.Pi)
	return T(int64(math.Floor(m + 0.5)))
}

func (f T) MarhsalJSON() ([]byte, error) {
	return json.Marshal(int64(f / Hertz))
}

func (f *T) UnmarshalJSON(b []byte) error {
	var i int64
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	*f = T(i) * Hertz
	return nil
}
