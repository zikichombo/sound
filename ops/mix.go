// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package ops

import (
	"errors"

	"github.com/irifrance/snd"
)

func MixEven(srcs ...snd.Source) (snd.Source, error) {
	n := len(srcs)
	if n == 0 {
		return nil, errors.New("cannot mix 0 sources.")
	}
	added, err := Add(srcs...)
	if err != nil {
		return nil, err
	}
	// we can amplify here since we are using float64 in memory
	return Amplify(added, 1.0/float64(n)), nil
}

func MustMixEven(srcs ...snd.Source) snd.Source {
	res, err := MixEven(srcs...)
	if err != nil {
		panic(err.Error())
	}
	return res
}
