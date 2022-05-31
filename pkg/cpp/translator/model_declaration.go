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
	"go/token"
)

type CTypeKind int

const (
	TypeKind CTypeKind = iota
	PlainTypeKind
	StructKind
	OpaqueStructKind
	UnionKind
	FunctionKind
	EnumKind
)

type CType interface {
	GetBase() string
	GetTag() string
	SetRaw(x string)
	CGoName() string
	GetPointers() uint8
	SetPointers(uint8)
	AddOuterArr(uint64)
	AddInnerArr(uint64)
	OuterArrays() ArraySpec
	InnerArrays() ArraySpec
	OuterArraySizes() []ArraySizeSpec
	InnerArraySizes() []ArraySizeSpec
	//
	IsConst() bool
	IsOpaque() bool
	IsComplete() bool
	Kind() CTypeKind
	String() string
	Copy() CType
	AtLevel(level int) CType
}

type (
	Value interface{}
)

type CDecl struct {
	Spec       CType
	Name       string
	Value      Value
	Expression string
	IsStatic   bool
	IsTypedef  bool
	IsDefine   bool
	Pos        token.Pos
	Src        string
}

func (c CDecl) String() string {
	buf := new(bytes.Buffer)
	switch {
	case len(c.Name) > 0:
		fmt.Fprintf(buf, "%s %s", c.Spec, c.Name)
	default:
		buf.WriteString(c.Spec.String())
	}
	if len(c.Expression) > 0 {
		fmt.Fprintf(buf, " = %s", string(c.Expression))
	}
	return buf.String()
}
