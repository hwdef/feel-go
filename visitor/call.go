package visitor

import (
	"github.com/kulics/lite-go/parser"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

func (sf *LiteVisitor) VisitCallExpression(ctx *parser.CallExpressionContext) interface{} {
	count := ctx.GetChildCount()
	r := Result{}
	if count == 3 {
		e1 := sf.Visit(ctx.GetChild(0).(antlr.ParseTree)).(Result)
		op := sf.Visit(ctx.GetChild(1).(antlr.ParseTree))
		e2 := sf.Visit(ctx.GetChild(2).(antlr.ParseTree)).(Result)
		r.Text = e1.Text + op.(string) + e2.Text
	} else if count == 1 {
		r = sf.Visit(ctx.GetChild(0).(antlr.ParseTree)).(Result)
	}
	return r
}

func (sf *LiteVisitor) VisitCallChannel(ctx *parser.CallChannelContext) interface{} {
	id := sf.Visit(ctx.Id()).(Result)
	if ctx.GetOp() != nil && ctx.GetOp().GetTokenType() == parser.LiteLexerQuestion {
		id.Text += "?"
	}
	id.Text = "<-" + id.Text
	return id
}

func (sf *LiteVisitor) VisitCallElement(ctx *parser.CallElementContext) interface{} {
	id := sf.Visit(ctx.Id()).(Result)
	if ctx.GetOp() != nil && ctx.GetOp().GetTokenType() == parser.LiteLexerQuestion {
		id.Text += "?"
	}
	if ctx.Expression() == nil {
		return Result{Text: id.Text + sf.Visit(ctx.Slice()).(string)}
	}
	r := sf.Visit(ctx.Expression()).(Result)
	r.Text = id.Text + "[" + r.Text + "]"
	return r
}

func (sf *LiteVisitor) VisitCallFunc(ctx *parser.CallFuncContext) interface{} {
	r := Result{Data: "var"}
	id := sf.Visit(ctx.Id()).(Result)
	r.Text += id.Text
	if ctx.TemplateCall() != nil {
		r.Text += sf.Visit(ctx.TemplateCall()).(string)
	}
	if ctx.Tuple() != nil {
		r.Text += "(" + sf.Visit(ctx.Tuple()).(Result).Text + ")"
	} else {
		r.Text += "(" + sf.Visit(ctx.Lambda()).(Result).Text + ")"
	}
	return r
}

func (sf *LiteVisitor) VisitCallNew(ctx *parser.CallNewContext) interface{} {
	r := Result{Data: sf.Visit(ctx.TypeType())}
	param := ""
	if ctx.ExpressionList() != nil {
		param = sf.Visit(ctx.ExpressionList()).(Result).Text
	}
	r.Text = "make(" + sf.Visit(ctx.TypeType()).(string)
	if param != "" {
		r.Text += "," + param
	}
	r.Text += ")"
	return r
}

func (sf *LiteVisitor) VisitCallPkg(ctx *parser.CallPkgContext) interface{} {
	r := Result{Data: sf.Visit(ctx.TypeType())}
	r.Text = sf.Visit(ctx.TypeType()).(string)
	if ctx.PkgAssign() != nil {
		r.Text += sf.Visit(ctx.PkgAssign()).(string)
	} else if ctx.ListAssign() != nil {
		r.Text += sf.Visit(ctx.ListAssign()).(string)
	} else if ctx.SetAssign() != nil {
		r.Text += sf.Visit(ctx.SetAssign()).(string)
	} else if ctx.DictionaryAssign() != nil {
		r.Text += sf.Visit(ctx.DictionaryAssign()).(string)
	} else {
		r.Text += "{}"
	}
	return r
}

func (sf *LiteVisitor) VisitPkgAssign(ctx *parser.PkgAssignContext) interface{} {
	obj := ""
	obj += "{"
	for i := 0; i < len(ctx.AllPkgAssignElement()); i++ {
		if i == 0 {
			obj += sf.Visit(ctx.PkgAssignElement(i)).(string)
		} else {
			obj += "," + sf.Visit(ctx.PkgAssignElement(i)).(string)
		}
	}
	obj += "}"
	return obj
}

func (sf *LiteVisitor) VisitListAssign(ctx *parser.ListAssignContext) interface{} {
	obj := ""
	obj += "{"
	for i := 0; i < len(ctx.AllExpression()); i++ {
		r := sf.Visit(ctx.Expression(i)).(Result)
		if i == 0 {
			obj += r.Text
		} else {
			obj += "," + r.Text
		}
	}
	obj += "}"
	return obj
}

func (sf *LiteVisitor) VisitSetAssign(ctx *parser.SetAssignContext) interface{} {
	obj := ""
	obj += "{"
	for i := 0; i < len(ctx.AllExpression()); i++ {
		r := sf.Visit(ctx.Expression(i)).(Result)
		if i == 0 {
			obj += r.Text
		} else {
			obj += "," + r.Text
		}
	}
	obj += "}"
	return obj
}

func (sf *LiteVisitor) VisitDictionaryAssign(ctx *parser.DictionaryAssignContext) interface{} {
	obj := ""
	obj += "{"
	for i := 0; i < len(ctx.AllDictionaryElement()); i++ {
		r := sf.Visit(ctx.DictionaryElement(i)).(DicEle)
		if i == 0 {
			obj += r.Text
		} else {
			obj += "," + r.Text
		}
	}
	obj += "}"
	return obj
}

func (sf *LiteVisitor) VisitPkgAssignElement(ctx *parser.PkgAssignElementContext) interface{} {
	obj := ""
	obj += sf.Visit(ctx.Name()).(string) + " = " + sf.Visit(ctx.Expression()).(Result).Text
	return obj
}

func (sf *LiteVisitor) VisitDictionaryElement(ctx *parser.DictionaryElementContext) interface{} {
	r1 := sf.Visit(ctx.Expression(0)).(Result)
	r2 := sf.Visit(ctx.Expression(1)).(Result)
	r := DicEle{
		Key:   r1.Data.(string),
		Value: r2.Data.(string),
		Text:  r1.Text + ":" + r2.Text,
	}
	return r
}

type DicEle struct {
	Key   string
	Value string
	Text  string
}
