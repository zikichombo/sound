// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package sample

import "math"

// ToDb translates v to decibels.
func ToDb(v float64) float64 {
	return 20 * math.Log10(v)
}

// FromDb translates db to a standard value.
func FromDb(db float64) float64 {
	db /= 20
	return math.Pow(10, db)
}
