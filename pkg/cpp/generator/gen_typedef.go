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
	"path/filepath"

	tl "github.com/bhojpur/build/pkg/cpp/translator"
)

func (gen *Generator) writeTypeTypedef(wr io.Writer, decl *tl.CDecl, seenNames map[string]bool) {
	goSpec := gen.tr.TranslateSpec(decl.Spec)
	goTypeName := gen.tr.TransformName(tl.TargetType, decl.Name)
	if seenNames[string(goTypeName)] {
		return
	} else {
		seenNames[string(goTypeName)] = true
	}
	fmt.Fprintf(wr, "// %s type as declared in %s\n", goTypeName,
		filepath.ToSlash(gen.tr.SrcLocation(tl.TargetType, decl.Name, decl.Pos)))
	fmt.Fprintf(wr, "type %s %s", goTypeName, goSpec.UnderlyingString())
	writeSpace(wr, 1)
}

func (gen *Generator) writeEnumTypedef(wr io.Writer, decl *tl.CDecl) {
	cName, ok := getName(decl)
	if !ok {
		return
	}
	goName := gen.tr.TransformName(tl.TargetType, cName)
	typeRef := gen.tr.TranslateSpec(decl.Spec).UnderlyingString()
	if typeName := string(goName); typeName != typeRef {
		fmt.Fprintf(wr, "// %s as declared in %s\n", goName,
			filepath.ToSlash(gen.tr.SrcLocation(tl.TargetConst, cName, decl.Pos)))
		fmt.Fprintf(wr, "type %s %s", goName, typeRef)
		writeSpace(wr, 1)
	}
}

func (gen *Generator) writeFunctionTypedef(wr io.Writer, decl *tl.CDecl, seenNames map[string]bool) {
	var returnRef string
	funcSpec := decl.Spec.Copy().(*tl.CFunctionSpec)
	funcSpec.Pointers = 0 // function pointers not supported here

	if funcSpec.Return != nil {
		// defaults to ref for the return values
		ptrTip := tl.TipPtrRef
		if ptrTipRx, ok := gen.tr.PtrTipRx(tl.TipScopeFunction, decl.Name); ok {
			if tip := ptrTipRx.Self(); tip.IsValid() {
				ptrTip = tip
			}
		}
		typeTip := tl.TipTypeNamed
		if typeTipRx, ok := gen.tr.TypeTipRx(tl.TipScopeFunction, decl.Name); ok {
			if tip := typeTipRx.Self(); tip.IsValid() {
				typeTip = tip
			}
		}
		returnRef = gen.tr.TranslateSpec(funcSpec.Return, ptrTip, typeTip).String()
	}

	ptrTipRx, _ := gen.tr.PtrTipRx(tl.TipScopeFunction, decl.Name)
	typeTipRx, _ := gen.tr.TypeTipRx(tl.TipScopeFunction, decl.Name)
	goFuncName := gen.tr.TransformName(tl.TargetType, decl.Name)
	if seenNames[string(goFuncName)] {
		return
	} else {
		seenNames[string(goFuncName)] = true
	}
	goSpec := gen.tr.TranslateSpec(funcSpec, ptrTipRx.Self(), typeTipRx.Self())
	goSpec.Raw = "" // not used in func typedef
	fmt.Fprintf(wr, "// %s type as declared in %s\n", goFuncName,
		filepath.ToSlash(gen.tr.SrcLocation(tl.TargetFunction, decl.Name, decl.Pos)))
	fmt.Fprintf(wr, "type %s %s", goFuncName, goSpec)
	gen.writeFunctionParams(wr, decl.Name, decl.Spec)
	if len(returnRef) > 0 {
		fmt.Fprintf(wr, " %s", returnRef)
	}
	for _, helper := range gen.getCallbackHelpers(string(goFuncName), decl.Name, decl.Spec) {
		gen.submitHelper(helper)
	}
	writeSpace(wr, 1)
}

func getName(decl *tl.CDecl) (string, bool) {
	if len(decl.Name) > 0 {
		return decl.Name, true
	}
	if base := decl.Spec.GetBase(); len(base) > 0 {
		return base, true
	}
	return "", false
}

func (gen *Generator) writeStructTypedef(wr io.Writer, decl *tl.CDecl, raw bool, seenNames map[string]bool) {
	cName, ok := getName(decl)
	if !ok {
		return
	}
	goName := gen.tr.TransformName(tl.TargetType, cName)
	if seenNames[string(goName)] {
		return
	} else {
		seenNames[string(goName)] = true
	}
	if raw || !decl.Spec.IsComplete() {
		// opaque struct
		fmt.Fprintf(wr, "// %s as declared in %s\n", goName,
			filepath.ToSlash(gen.tr.SrcLocation(tl.TargetType, cName, decl.Pos)))
		fmt.Fprintf(wr, "type %s C.%s", goName, decl.Spec.CGoName())
		writeSpace(wr, 1)
		for _, helper := range gen.getRawStructHelpers(goName, cName, decl.Spec) {
			gen.submitHelper(helper)
		}
		return
	}

	fmt.Fprintf(wr, "// %s as declared in %s\n", goName,
		filepath.ToSlash(gen.tr.SrcLocation(tl.TargetType, cName, decl.Pos)))
	fmt.Fprintf(wr, "type %s struct {", goName)
	writeSpace(wr, 1)
	gen.submitHelper(cgoAllocMap)
	gen.writeStructMembers(wr, cName, decl.Spec)
	writeEndStruct(wr)
	writeSpace(wr, 1)
	for _, helper := range gen.getStructHelpers(goName, cName, decl.Spec) {
		gen.submitHelper(helper)
	}
}

func (gen *Generator) writeUnionTypedef(wr io.Writer, decl *tl.CDecl) {
	cName, ok := getName(decl)
	if !ok {
		return
	}
	goName := gen.tr.TransformName(tl.TargetType, cName)
	typeRef := gen.tr.TranslateSpec(decl.Spec).UnderlyingString()

	if typeName := string(goName); typeName != typeRef {
		fmt.Fprintf(wr, "// %s as declared in %s\n", goName,
			filepath.ToSlash(gen.tr.SrcLocation(tl.TargetType, cName, decl.Pos)))
		fmt.Fprintf(wr, "const sizeof%s = unsafe.Sizeof(C.%s{})\n", goName, decl.Spec.CGoName())
		fmt.Fprintf(wr, "type %s [sizeof%s]byte\n", goName, goName)
		writeSpace(wr, 1)
		return
	}
}
