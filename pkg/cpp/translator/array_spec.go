package translator

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"
	"strconv"
	"strings"
)

type ArraySpec string

func (a ArraySpec) String() string {
	return string(a)
}

func (a *ArraySpec) AddSized(size uint64) {
	*a = ArraySpec(fmt.Sprintf("%s[%d]", a, size))
}

func (a *ArraySpec) Prepend(spec ArraySpec) {
	*a = spec + *a
}

type ArraySizeSpec struct {
	N   uint64
	Str string
}

func (a ArraySpec) Sizes() (sizes []ArraySizeSpec) {
	if len(a) == 0 {
		return
	}
	arr := string(a)
	for len(arr) > 0 {
		// get "n" from "[k][l][m][n]"
		p1 := strings.LastIndexByte(arr, '[')
		p2 := strings.LastIndexByte(arr, ']')
		part := arr[p1+1 : p2]
		// and try to convert uint64
		if u, err := strconv.ParseUint(part, 10, 64); err != nil || u == 0 {
			// use size spec as-is (i.e. unsafe.Sizeof(x))
			sizes = append(sizes, ArraySizeSpec{Str: part})
		} else {
			sizes = append(sizes, ArraySizeSpec{N: u})
		}
		arr = arr[:p1]
	}
	return sizes
}
