package builder

import (
	"fmt"

	"github.com/elz-lang/elz/src/elz/ast"
	"github.com/elz-lang/elz/src/elz/parser"
)

type BindingBuilder struct {
	bindingType  map[string][]ast.Type
	bindTypeList []ast.Type

	binding *ast.Binding
}

func NewBindingBuilder() *BindingBuilder {
	return &BindingBuilder{
		bindingType:  make(map[string][]ast.Type),
		bindTypeList: make([]ast.Type, 0),
	}
}

func (b *BindingBuilder) ExitBindType(c *parser.BindTypeContext) {
	bindName := c.IDENT().GetText()
	_, exist := b.bindingType[bindName]
	if exist {
		panic("bind type existed")
	}
	b.bindingType[bindName] = b.bindTypeList
}

// ExitExistType listen format: `int` represents existing type
// You would get compile fatal if the type isn't existed
func (b *BindingBuilder) ExitExistType(c *parser.ExistTypeContext) {
	b.bindTypeList = append(b.bindTypeList, &ast.ExistType{Name: c.IDENT().GetText()})
}

func (b *BindingBuilder) ExitVoidType(c *parser.VoidTypeContext) {
	b.bindTypeList = append(b.bindTypeList, &ast.VoidType{})
}

// ExitVariantType listen format: `'a` as type hole
func (b *BindingBuilder) ExitVariantType(c *parser.VariantTypeContext) {
	b.bindTypeList = append(b.bindTypeList, &ast.VariantType{Name: c.IDENT().GetText()})
}

// ExitCombineType listen format: `int -> int -> int`
func (b *BindingBuilder) ExitCombineType(c *parser.CombineTypeContext) {
	// ignore, just help we know has this syntax
}

func (b *BindingBuilder) Build(c *parser.BindingContext, bindingTo ast.Expr) *ast.Binding {
	paramList := make([]string, 0)
	for _, paramName := range c.AllIDENT() {
		paramList = append(paramList, paramName.GetText())
	}
	bindName := c.IDENT(0).GetText()
	export := c.KEYWORD_EXPORT() != nil
	binding := &ast.Binding{
		Export:    export,
		Name:      bindName,
		ParamList: paramList[1:],
		Expr:      bindingTo,
	}
	if t, exist := b.bindingType[bindName]; exist {
		binding.Type = t
	}
	return binding
}

func (b *Builder) ExitBinding(c *parser.BindingContext) {
	binding := b.BindingBuilder.Build(c, b.ExprBuilder.Build())
	err := b.astTree.InsertBinding(binding)
	if err != nil {
		err := fmt.Errorf("stop parsing, error: %s", err)
		panic(err)
	}
}
