package freq

import (
	"fmt"
	"math"
)

// Mel is a unit of frequency in which a difference less than MelAudible
// is generally imperceivable (melodically).
type Mel int64

// MelAudible is the smallest difference in frequency in Mels which
// is melodically audible.  Smaller harmonic differences are audible
// due to beating.
const MelAudible Mel = Mel(1000000000)

// ToMel gives f in Mels.
func ToMel(f T) Mel {
	sign := int64(1)
	if f < 0 {
		f = -f
		sign = -1
	}
	hz := f.Float64()
	v := 1.0 + hz/700.0
	v = 2595 * math.Log10(v)
	m := sign * int64(math.Floor(v+0.5)) * int64(MelAudible)
	return Mel(m)
}

// String returns a string representation of m.
func (m Mel) String() string {
	return fmt.Sprintf("%d.%03dmel", m/MelAudible, (m%MelAudible)/Mel(MilliHertz))
}

// FromMel gives m as frequency.
func FromMel(m Mel) T {
	return m.Freq()
}

// Freq gives m as frequency.
func (m Mel) Freq() T {
	sign := int64(1)
	if m < 0 {
		m = -m
		sign = -1
	}
	v := float64(int64(m) / int64(Hertz))
	v /= 2595
	v = math.Pow(10, v)
	v -= 1
	v *= 700
	f := sign * int64(math.Floor(v+0.5)) * int64(Hertz)
	return T(f)
}
