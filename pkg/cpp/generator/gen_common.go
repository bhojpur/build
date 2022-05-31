package generator

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
	"io"
	"strings"

	tl "github.com/bhojpur/build/pkg/cpp/translator"
)

var (
	skipName    = []byte("_")
	skipNameStr = "_"
)

func (gen *Generator) writeStructMembers(wr io.Writer, structName string, spec tl.CType) {
	structSpec := spec.(*tl.CStructSpec)
	ptrTipRx, typeTipRx, memTipRx := gen.tr.TipRxsForSpec(tl.TipScopeType, structName, structSpec)
	const public = true
	for i, member := range structSpec.Members {
		ptrTip := ptrTipRx.TipAt(i)
		if !ptrTip.IsValid() {
			ptrTip = tl.TipPtrArr
		}
		typeTip := typeTipRx.TipAt(i)
		if !typeTip.IsValid() {
			typeTip = tl.TipTypeNamed
		}
		memTip := memTipRx.TipAt(i)
		if !memTip.IsValid() {
			memTip = gen.MemTipOf(member)
		}
		if memTip == tl.TipMemRaw {
			ptrTip = tl.TipPtrSRef
		}
		declName := checkName(gen.tr.TransformName(tl.TargetType, member.Name, public))
		switch member.Spec.Kind() {
		case tl.TypeKind:
			goSpec := gen.tr.TranslateSpec(member.Spec, ptrTip, typeTip)
			fmt.Fprintf(wr, "%s %s", declName, goSpec)
		case tl.StructKind, tl.OpaqueStructKind, tl.UnionKind:
			if !gen.tr.IsAcceptableName(tl.TargetType, member.Spec.GetBase()) {
				continue
			}
			goSpec := gen.tr.TranslateSpec(member.Spec, ptrTip, typeTip)
			fmt.Fprintf(wr, "%s %s", declName, goSpec)
		case tl.EnumKind:
			if !gen.tr.IsAcceptableName(tl.TargetType, member.Spec.GetBase()) {
				continue
			}
			typeRef := gen.tr.TranslateSpec(member.Spec, ptrTip, typeTip).String()
			fmt.Fprintf(wr, "%s %s", declName, typeRef)
		case tl.FunctionKind:
			gen.writeFunctionAsArg(wr, member, ptrTip, typeTip, public)
		}
		writeSpace(wr, 1)
	}

	if memTipRx.Self() == tl.TipMemRaw {
		return
	}

	crc := getRefCRC(structSpec)
	cgoSpec := gen.tr.CGoSpec(structSpec, false)
	if len(cgoSpec.Base) == 0 {
		return
	}
	fmt.Fprintf(wr, "ref%2x *%s\n", crc, cgoSpec)
	fmt.Fprintf(wr, "allocs%2x interface{}\n", crc)
}

func (gen *Generator) writeInstanceObjectParam(wr io.Writer, funcName string, funcSpec tl.CType) {
	spec := funcSpec.(*tl.CFunctionSpec)
	ptrTipSpecRx, _ := gen.tr.PtrTipRx(tl.TipScopeFunction, funcName)
	typeTipSpecRx, _ := gen.tr.TypeTipRx(tl.TipScopeFunction, funcName)

	for i, param := range spec.Params {
		ptrTip := ptrTipSpecRx.TipAt(i)

		if ptrTip != tl.TipPtrInst {
			continue
		}

		ptrTip = tl.TipPtrRef
		typeTip := typeTipSpecRx.TipAt(i)
		if !typeTip.IsValid() {
			// try to use type tip for the type itself
			if tip, ok := gen.tr.TypeTipRx(tl.TipScopeType, param.Spec.CGoName()); ok {
				if tip := tip.Self(); tip.IsValid() {
					typeTip = tip
				}
			}
		}

		writeSpace(wr, 1)
		writeStartParams(wr)
		gen.writeFunctionParam(wr, param, ptrTip, typeTip)
		writeEndParams(wr)

		break
	}
}

func (gen *Generator) writeFunctionParam(wr io.Writer, param *tl.CDecl, ptrTip tl.Tip, typeTip tl.Tip) {
	const public = false

	declName := checkName(gen.tr.TransformName(tl.TargetType, param.Name, public))
	switch param.Spec.Kind() {
	case tl.TypeKind:
		goSpec := gen.tr.TranslateSpec(param.Spec, ptrTip, typeTip)
		if len(goSpec.OuterArr) > 0 {
			fmt.Fprintf(wr, "%s *%s", declName, goSpec)
		} else {
			fmt.Fprintf(wr, "%s %s", declName, goSpec)
		}
	case tl.StructKind, tl.OpaqueStructKind, tl.UnionKind:
		goSpec := gen.tr.TranslateSpec(param.Spec, ptrTip, typeTip)
		if len(goSpec.OuterArr) > 0 {
			fmt.Fprintf(wr, "%s *%s", declName, goSpec)
		} else {
			fmt.Fprintf(wr, "%s %s", declName, goSpec)
		}
	case tl.EnumKind:
		typeRef := gen.tr.TranslateSpec(param.Spec, ptrTip, typeTip).String()
		fmt.Fprintf(wr, "%s %s", declName, typeRef)
	case tl.FunctionKind:
		gen.writeFunctionAsArg(wr, param, ptrTip, typeTip, public)
	}
}

func (gen *Generator) writeFunctionParams(wr io.Writer, funcName string, funcSpec tl.CType) {
	spec := funcSpec.(*tl.CFunctionSpec)
	ptrTipSpecRx, _ := gen.tr.PtrTipRx(tl.TipScopeFunction, funcName)
	typeTipSpecRx, _ := gen.tr.TypeTipRx(tl.TipScopeFunction, funcName)

	writeStartParams(wr)
	for i, param := range spec.Params {
		ptrTip := ptrTipSpecRx.TipAt(i)

		if ptrTip == tl.TipPtrInst {
			continue
		}

		if !ptrTip.IsValid() {
			ptrTip = tl.TipPtrArr
		}

		typeTip := typeTipSpecRx.TipAt(i)
		if !typeTip.IsValid() {
			// try to use type tip for the type itself
			if tip, ok := gen.tr.TypeTipRx(tl.TipScopeType, param.Spec.CGoName()); ok {
				if tip := tip.Self(); tip.IsValid() {
					typeTip = tip
				}
			}
		}

		gen.writeFunctionParam(wr, param, ptrTip, typeTip)

		if i < len(spec.Params)-1 && ptrTipSpecRx.TipAt(i+1) != tl.TipPtrInst {
			fmt.Fprintf(wr, ", ")
		}
	}
	writeEndParams(wr)
}

func writeStartParams(wr io.Writer) {
	fmt.Fprint(wr, "(")
}

func writeEndParams(wr io.Writer) {
	fmt.Fprint(wr, ")")
}

func writeEndStruct(wr io.Writer) {
	fmt.Fprint(wr, "}")
}

func writeStartFuncBody(wr io.Writer) {
	fmt.Fprintln(wr, "{")
}

func writeEndFuncBody(wr io.Writer) {
	fmt.Fprintln(wr, "}")
}

func writeSpace(wr io.Writer, n int) {
	fmt.Fprint(wr, strings.Repeat("\n", n))
}

func writeError(wr io.Writer, err error) {
	fmt.Fprintf(wr, "// error: %v\n", err)
}
