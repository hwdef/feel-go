package visitor

import "github.com/kulics/lite-go/parser"

type Iterator struct {
	Begin  Result
	End    Result
	Step   Result
	Order  bool
	Attach bool
}

func (sf *LiteVisitor) VisitIteratorStatement(ctx *parser.IteratorStatementContext) interface{} {
	it := Iterator{Order: true, Attach: false}
	if ctx.GetOp().GetText() == ">=" || ctx.GetOp().GetText() == "<=" {
		it.Attach = true
	}
	if ctx.GetOp().GetText() == ">" || ctx.GetOp().GetText() == ">=" {
		it.Order = false
	}
	if len(ctx.AllExpression()) == 2 {
		it.Begin = sf.Visit(ctx.Expression(0)).(Result)
		it.End = sf.Visit(ctx.Expression(1)).(Result)
		it.Step = Result{Data: I32, Text: "1"}
	} else {
		it.Begin = sf.Visit(ctx.Expression(0)).(Result)
		it.End = sf.Visit(ctx.Expression(1)).(Result)
		it.Step = sf.Visit(ctx.Expression(2)).(Result)
	}
	return it
}

func (sf *LiteVisitor) VisitLoopStatement(ctx *parser.LoopStatementContext) interface{} {
	obj := ""
	id := "ea"
	if ctx.Id() != nil {
		id = sf.Visit(ctx.Id()).(Result).Text
	}
	it := sf.Visit(ctx.IteratorStatement()).(Iterator)
	order := ""
	step := ""
	if it.Order {
		step = "+="
		if it.Attach {
			order = "<="
		} else {
			order = "<"
		}
	} else {
		step = "-="
		if it.Attach {
			order = ">="
		} else {
			order = ">"
		}
	}
	obj += "for " + id + " := " + it.Begin.Text + ";" + id + order + it.End.Text + ";" + id + step + it.Step.Text

	obj += BlockLeft + Wrap
	obj += sf.ProcessFunctionSupport(ctx.AllFunctionSupportStatement())
	obj += BlockRight + Wrap
	return obj
}

func (sf *LiteVisitor) VisitLoopInfiniteStatement(ctx *parser.LoopInfiniteStatementContext) interface{} {
	obj := "for " + BlockLeft + Wrap
	obj += sf.ProcessFunctionSupport(ctx.AllFunctionSupportStatement())
	obj += BlockRight + Wrap
	return obj
}

func (sf *LiteVisitor) VisitLoopEachStatement(ctx *parser.LoopEachStatementContext) interface{} {
	obj := ""
	arr := sf.Visit(ctx.Expression()).(Result)
	target := arr.Text
	id := "ea"
	if len(ctx.AllId()) == 2 {
		id = sf.Visit(ctx.Id(0)).(Result).Text + "," + sf.Visit(ctx.Id(1)).(Result).Text
	} else if len(ctx.AllId()) == 1 {
		id = "_," + sf.Visit(ctx.Id(0)).(Result).Text
	}

	obj += "for " + id + " := range " + target
	obj += BlockLeft + Wrap
	obj += sf.ProcessFunctionSupport(ctx.AllFunctionSupportStatement())
	obj += BlockRight + Wrap
	return obj
}

func (sf *LiteVisitor) VisitLoopCaseStatement(ctx *parser.LoopCaseStatementContext) interface{} {
	obj := ""
	expr := sf.Visit(ctx.Expression()).(Result)
	obj += "for " + expr.Text
	obj += BlockLeft + Wrap
	obj += sf.ProcessFunctionSupport(ctx.AllFunctionSupportStatement())
	obj += BlockRight + Wrap
	return obj
}

func (sf *LiteVisitor) VisitLoopJumpStatement(ctx *parser.LoopJumpStatementContext) interface{} {
	return "break" + Wrap
}

func (sf *LiteVisitor) VisitLoopContinueStatement(ctx *parser.LoopContinueStatementContext) interface{} {
	return "continue" + Wrap
}
