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
	"strings"
)

type CFunctionSpec struct {
	Raw      string
	Typedef  string
	Return   CType
	Params   []*CDecl
	Pointers uint8
}

func (c CFunctionSpec) String() string {
	if len(c.Raw) > 0 {
		return c.Raw
	}
	var params []string
	for i, param := range c.Params {
		if len(param.Name) == 0 {
			params = append(params, fmt.Sprintf("arg%d", i))
			continue
		}
		params = append(params, param.String())
	}
	paramList := strings.Join(params, ", ")
	if c.Return != nil {
		return fmt.Sprintf("%s (%s)", c.Return, paramList)
	}
	return fmt.Sprintf("void (%s)", paramList)
}

func (c *CFunctionSpec) SetPointers(n uint8) {
	c.Pointers = n
}

func (c *CFunctionSpec) Kind() CTypeKind {
	return FunctionKind
}

func (c *CFunctionSpec) IsComplete() bool {
	return true
}

func (c *CFunctionSpec) IsOpaque() bool {
	return len(c.Params) == 0
}

func (c CFunctionSpec) Copy() CType {
	return &c
}

func (c *CFunctionSpec) GetBase() string {
	return ""
}

func (c *CFunctionSpec) GetTag() string {
	return c.Raw
}

func (c *CFunctionSpec) CGoName() string {
	return c.Raw
}

func (c *CFunctionSpec) SetRaw(x string) {
	c.Raw = x
}

func (c *CFunctionSpec) AddOuterArr(uint64) {}
func (c *CFunctionSpec) AddInnerArr(uint64) {}
func (c *CFunctionSpec) OuterArraySizes() []ArraySizeSpec {
	return nil
}
func (c *CFunctionSpec) InnerArraySizes() []ArraySizeSpec {
	return nil
}
func (c *CFunctionSpec) OuterArrays() ArraySpec {
	return ""
}
func (c *CFunctionSpec) InnerArrays() ArraySpec {
	return ""
}

func (c *CFunctionSpec) GetPointers() uint8 {
	return c.Pointers
}

func (c *CFunctionSpec) IsConst() bool {
	return false
}

func (c CFunctionSpec) AtLevel(level int) CType {
	return &c
}
