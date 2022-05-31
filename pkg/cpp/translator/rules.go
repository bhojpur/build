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

type Rules map[RuleTarget][]RuleSpec
type ConstRules map[ConstScope]ConstRule
type PtrTips map[TipScope][]TipSpec
type TypeTips map[TipScope][]TipSpec
type MemTips []TipSpec

type RuleSpec struct {
	From, To  string
	Action    RuleAction
	Transform RuleTransform
	Load      string
}

func (r *RuleSpec) LoadSpec(r2 RuleSpec) {
	if len(r.From) == 0 {
		r.From = r2.From
	}
	if len(r.To) == 0 {
		r.To = r2.To
	}
	if len(r.Action) == 0 {
		r.Action = r2.Action
	}
	if len(r.Transform) == 0 {
		r.Transform = r2.Transform
	}
}

type RuleAction string

const (
	ActionNone     RuleAction = ""
	ActionAccept   RuleAction = "accept"
	ActionIgnore   RuleAction = "ignore"
	ActionReplace  RuleAction = "replace"
	ActionDocument RuleAction = "doc"
)

var ruleActions = []RuleAction{
	ActionAccept, ActionIgnore, ActionReplace, ActionDocument,
}

type RuleTransform string

const (
	TransformLower    RuleTransform = "lower"
	TransformTitle    RuleTransform = "title"
	TransformExport   RuleTransform = "export"
	TransformUnexport RuleTransform = "unexport"
	TransformUpper    RuleTransform = "upper"
)

type RuleTarget string

const (
	NoTarget         RuleTarget = ""
	TargetGlobal     RuleTarget = "global"
	TargetPostGlobal RuleTarget = "post-global"
	//
	TargetConst    RuleTarget = "const"
	TargetType     RuleTarget = "type"
	TargetFunction RuleTarget = "function"
	//
	TargetPublic  RuleTarget = "public"
	TargetPrivate RuleTarget = "private"
)

type ConstRule string

const (
	ConstCGOAlias ConstRule = "cgo"
	ConstExpand   ConstRule = "expand"
	ConstEval     ConstRule = "eval"
)

type ConstScope string

const (
	ConstEnum    ConstScope = "enum"
	ConstDecl    ConstScope = "decl"
	ConstDefines ConstScope = "defines"
)

type Tip string

const (
	TipPtrSRef    Tip = "sref"
	TipPtrRef     Tip = "ref"
	TipPtrArr     Tip = "arr"
	TipPtrInst    Tip = "inst"
	TipMemRaw     Tip = "raw"
	TipTypeNamed  Tip = "named"
	TipTypePlain  Tip = "plain"
	TipTypeString Tip = "string"
	NoTip         Tip = ""
)

type TipKind string

const (
	TipKindUnknown TipKind = "unknown"
	TipKindPtr     TipKind = "ptr"
	TipKindType    TipKind = "type"
	TipKindMem     TipKind = "mem"
)

func (t Tip) Kind() TipKind {
	switch t {
	case TipPtrArr, TipPtrRef, TipPtrSRef, TipPtrInst:
		return TipKindPtr
	case TipTypePlain, TipTypeNamed, TipTypeString:
		return TipKindType
	case TipMemRaw:
		return TipKindMem
	default:
		return TipKindUnknown
	}
}

func (t Tip) IsValid() bool {
	switch t {
	case TipPtrArr, TipPtrRef, TipPtrSRef, TipPtrInst:
		return true
	case TipTypePlain, TipTypeNamed, TipTypeString:
		return true
	case TipMemRaw:
		return true
	default:
		return false
	}
}

type TipSpec struct {
	Target  string
	Tips    Tips
	Self    Tip
	Default Tip
}

type TipScope string

const (
	TipScopeAny      TipScope = "any"
	TipScopeStruct   TipScope = "struct"
	TipScopeType     TipScope = "type"
	TipScopeFunction TipScope = "function"
)

type Tips []Tip

var builtinRules = map[string]RuleSpec{
	"snakecase":  RuleSpec{Action: ActionReplace, From: "_([^_]+)", To: "$1", Transform: TransformTitle},
	"doc.file":   RuleSpec{Action: ActionDocument, To: "$path:$line"},
	"doc.google": RuleSpec{Action: ActionDocument, To: "https://google.com/search?q=$file+$name"},
}
