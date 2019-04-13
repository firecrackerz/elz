package builder

import (
	"fmt"
	"strconv"

	"github.com/elz-lang/elz/src/elz/ast"
	"github.com/elz-lang/elz/src/elz/internal/collection/stack"
	"github.com/elz-lang/elz/src/elz/parser"
)

type ExprBuilder struct {
	*parser.BaseElzListener
	exprStack *stack.Stack
}

func NewExprBuilder() *ExprBuilder {
	return &ExprBuilder{
		exprStack: stack.New(),
	}
}

func (b *ExprBuilder) Build() ast.Expr {
	return b.PopExpr().(ast.Expr)
}

func (b *ExprBuilder) PushExpr(e interface{}) {
	b.exprStack.Push(e)
}
func (b *ExprBuilder) PopExpr() ast.Expr {
	e := b.exprStack.Pop()
	if e != nil {
		return e.(ast.Expr)
	}
	return nil
}

// ExitMulDiv:
//
//  1 * 2
//  3 / 4
//
// WorkFlow:
//
// e.g. leftExpr Op rightExpr
// 1. push leftExpr
// 2. push rightExpr
// MulDiv
// 3. pop rightExpr
// 4. pop leftExpr
func (b *ExprBuilder) ExitMulDiv(c *parser.MulDivContext) {
	rExpr := b.PopExpr()
	lExpr := b.PopExpr()
	// left expr, right expr, operator
	b.PushExpr(&ast.BinaryExpr{
		LExpr: lExpr,
		RExpr: rExpr,
		Op:    c.GetOp().GetText(),
	})
}
func (b *ExprBuilder) ExitAddSub(c *parser.AddSubContext) {
	rExpr := b.PopExpr()
	lExpr := b.PopExpr()
	// left expr, right expr, operator
	b.PushExpr(&ast.BinaryExpr{
		LExpr: lExpr,
		RExpr: rExpr,
		Op:    c.GetOp().GetText(),
	})
}
func (b *ExprBuilder) ExitInt(c *parser.IntContext) {
	b.PushExpr(ast.NewInt(c.INT().GetText()))
}
func (b *ExprBuilder) ExitFloat(c *parser.FloatContext) {
	b.PushExpr(ast.NewFloat(c.FLOAT().GetText()))
}
func (b *ExprBuilder) ExitString(c *parser.StringContext) {
	literal := c.STRING().GetText()
	literal, err := strconv.Unquote(literal)
	if err != nil {
		panic(fmt.Errorf("failed at unquoting string literal, it's a compiler bug that you should report, error: %s", err))
	}
	b.PushExpr(ast.NewString(literal))
}
func (b *ExprBuilder) ExitBoolean(c *parser.BooleanContext) {
	b.PushExpr(ast.NewBool(c.BOOLEAN().GetText()))
}
func (b *ExprBuilder) ExitList(c *parser.ListContext) {
	exprList := make([]ast.Expr, 0)
	for e := b.PopExpr(); e != nil; e = b.PopExpr() {
		exprList = append(exprList, e)
	}
	b.PushExpr(ast.NewList(exprList...))
}
func (b *ExprBuilder) ExitIdentifier(c *parser.IdentifierContext) {
	b.PushExpr(ast.NewIdent(c.GetText()))
}

func (b *ExprBuilder) ExitFnCall(c *parser.FnCallContext) {
	s := stack.New()
	for i := 0; i < len(c.AllArg()); i++ {
		s.Push(b.PopExpr().(*ast.Arg))
	}
	exprList := make([]*ast.Arg, len(c.AllArg()))
	for i, _ := range exprList {
		e, isArg := s.Pop().(*ast.Arg)
		if isArg {
			exprList[i] = e
		} else {
			panic(fmt.Errorf("expression in function call is not an argument, it must be compiler bug, report it to the project, error: %#v", e))
		}
	}
	b.PushExpr(&ast.FuncCall{
		AccessChain: c.AccessChain().GetText(),
		ArgList:     exprList,
	})
}

func (b *ExprBuilder) ExitArg(c *parser.ArgContext) {
	expr := b.PopExpr()
	ident := ""
	if c.IDENT() != nil {
		ident = c.IDENT().GetText()
	}
	b.PushExpr(&ast.Arg{
		Ident: ident,
		Expr:  expr,
	})
}
