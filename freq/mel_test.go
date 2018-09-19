package freq_test

import (
	"math/rand"
	"testing"

	"zikichombo.org/sound/freq"
)

func TestMelConvert(t *testing.T) {
	N := 128
	for i := 0; i < N; i++ {
		f := freq.T(rand.Intn(int(20 * freq.KiloHertz)))
		m := freq.ToMel(f)
		g := m.Freq()
		i := freq.Diff(f, g)
		if i > 9*freq.Cent { // usually 0, but some cases have lots of rounding errors.
			t.Errorf("%s -> %s -> %s; d%s\n", f, m, g, i)
		}
	}
}
