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
	"bytes"
	"fmt"
)

type CTypeSpec struct {
	Raw      string
	Base     string
	Const    bool
	Signed   bool
	Unsigned bool
	Short    bool
	Long     bool
	Complex  bool
	Opaque   bool
	Pointers uint8
	InnerArr ArraySpec
	OuterArr ArraySpec
}

func (spec CTypeSpec) String() string {
	buf := new(bytes.Buffer)
	if spec.Unsigned {
		buf.WriteString("unsigned ")
	} else if spec.Signed {
		buf.WriteString("signed ")
	}
	switch {
	case spec.Long:
		buf.WriteString("long ")
	case spec.Short:
		buf.WriteString("short ")
	case spec.Complex:
		buf.WriteString("complex ")
	}
	fmt.Fprint(buf, spec.Base)

	var unsafePointer uint8
	if spec.Base == "unsafe.Pointer" {
		unsafePointer = 1
	}

	buf.WriteString(arrs(spec.InnerArr))
	buf.WriteString(ptrs(spec.Pointers - unsafePointer))
	buf.WriteString(arrs(spec.OuterArr))
	return buf.String()
}

func (c *CTypeSpec) SetPointers(n uint8) {
	c.Pointers = n
}

func (c *CTypeSpec) IsComplete() bool {
	return true
}

func (c *CTypeSpec) IsOpaque() bool {
	return c.Opaque
}

func (c *CTypeSpec) Kind() CTypeKind {
	return TypeKind
}

func (c CTypeSpec) Copy() CType {
	return &c
}

func (c *CTypeSpec) GetBase() string {
	return c.Base
}

func (c *CTypeSpec) GetTag() string {
	return ""
}

func (c *CTypeSpec) SetRaw(x string) {
	c.Raw = x
}

func (c *CTypeSpec) CGoName() (name string) {
	if len(c.Raw) > 0 {
		return c.Raw
	}
	switch c.Base {
	case "int", "short", "long", "char":
		if c.Unsigned {
			name += "u"
		} else if c.Signed {
			name += "s"
		}
		switch {
		case c.Long:
			name += "long"
			if c.Base == "long" {
				name += "long"
			}
		case c.Short:
			name += "short"
		default:
			name += c.Base
		}
	default:
		name = c.Base
	}
	return
}

func (c *CTypeSpec) AddOuterArr(size uint64) {
	c.OuterArr.AddSized(size)
}

func (c *CTypeSpec) AddInnerArr(size uint64) {
	c.InnerArr.AddSized(size)
}

func (c *CTypeSpec) OuterArraySizes() []ArraySizeSpec {
	return c.OuterArr.Sizes()
}

func (c *CTypeSpec) InnerArraySizes() []ArraySizeSpec {
	return c.InnerArr.Sizes()
}

func (c *CTypeSpec) OuterArrays() ArraySpec {
	return c.OuterArr
}

func (c *CTypeSpec) InnerArrays() ArraySpec {
	return c.InnerArr
}

func (c *CTypeSpec) GetPointers() uint8 {
	return c.Pointers
}

func (c *CTypeSpec) IsConst() bool {
	return c.Const
}

func (c CTypeSpec) AtLevel(level int) CType {
	spec := c
	var outerArrSpec ArraySpec
	for i, size := range spec.OuterArr.Sizes() {
		if i < int(level) {
			continue
		} else if i == 0 {
			spec.Pointers = 1
			continue
		}
		outerArrSpec.AddSized(size.N)
	}
	if int(level) > len(spec.OuterArr) {
		if delta := int(spec.Pointers) + len(spec.OuterArr.Sizes()) - int(level); delta > 0 {
			spec.Pointers = uint8(delta)
		}
	}
	spec.OuterArr = outerArrSpec
	spec.InnerArr = ArraySpec("")
	return &spec
}
