// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package freq

import (
	"fmt"
	"testing"
)

func TestTemper(t *testing.T) {
	st := Temper(12)
	c := Temper(1200)
	fmt.Printf("st %s\nc  %s\n", st, c)
}
