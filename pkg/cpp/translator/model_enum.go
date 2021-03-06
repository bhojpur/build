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
	"strings"
)

type CEnumSpec struct {
	Tag      string
	Typedef  string
	Members  []*CDecl
	Type     CTypeSpec
	Pointers uint8
	InnerArr ArraySpec
	OuterArr ArraySpec
}

func (c *CEnumSpec) PromoteType(v Value) *CTypeSpec {
	var (
		int32Spec = CTypeSpec{Base: "int"}
		int64Spec = CTypeSpec{Base: "long int"}
	)
	switch c.Type {
	case int32Spec: // need promotion
		switch v.(type) {
		case int64:
			c.Type = int64Spec
		}
	case int64Spec:
	default:
		switch v.(type) {
		case int32:
			c.Type = int32Spec
		case int64:
			c.Type = int64Spec
		}
	}
	return &c.Type
}

func (spec CEnumSpec) String() string {
	buf := new(bytes.Buffer)
	writePrefix := func() {
		buf.WriteString("enum ")
	}

	switch {
	case len(spec.Typedef) > 0:
		buf.WriteString(spec.Typedef)
	case len(spec.Tag) > 0:
		writePrefix()
		buf.WriteString(spec.Tag)
	case len(spec.Members) > 0:
		var members []string
		for _, m := range spec.Members {
			members = append(members, m.String())
		}
		membersColumn := strings.Join(members, ",\n")
		writePrefix()
		fmt.Fprintf(buf, " {%s}", membersColumn)
	default:
		writePrefix()
	}

	buf.WriteString(arrs(spec.OuterArr))
	buf.WriteString(ptrs(spec.Pointers))
	buf.WriteString(arrs(spec.InnerArr))
	return buf.String()
}

func (c *CEnumSpec) SetPointers(n uint8) {
	c.Pointers = n
}

func (c *CEnumSpec) Kind() CTypeKind {
	return EnumKind
}

func (c *CEnumSpec) IsComplete() bool {
	return len(c.Members) > 0
}

func (c *CEnumSpec) IsOpaque() bool {
	return len(c.Members) == 0
}

func (c CEnumSpec) Copy() CType {
	return &c
}

func (c *CEnumSpec) GetBase() string {
	if len(c.Typedef) > 0 {
		return c.Typedef
	}
	return c.Tag
}

func (c *CEnumSpec) SetRaw(x string) {
	c.Typedef = x
}

func (c *CEnumSpec) GetTag() string {
	return c.Tag
}

func (c *CEnumSpec) CGoName() string {
	if len(c.Typedef) > 0 {
		return c.Typedef
	}
	return "enum_" + c.Tag
}

func (c *CEnumSpec) AddOuterArr(size uint64) {
	c.OuterArr.AddSized(size)
}

func (c *CEnumSpec) AddInnerArr(size uint64) {
	c.InnerArr.AddSized(size)
}

func (c *CEnumSpec) OuterArraySizes() []ArraySizeSpec {
	return c.OuterArr.Sizes()
}

func (c *CEnumSpec) InnerArraySizes() []ArraySizeSpec {
	return c.InnerArr.Sizes()
}

func (c *CEnumSpec) OuterArrays() ArraySpec {
	return c.OuterArr
}

func (c *CEnumSpec) InnerArrays() ArraySpec {
	return c.InnerArr
}

func (c *CEnumSpec) GetPointers() uint8 {
	return c.Pointers
}

func (c *CEnumSpec) IsConst() bool {
	// could be c.Const
	return false
}

func (c CEnumSpec) AtLevel(level int) CType {
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
