package builder

import (
	"testing"

	"github.com/elz-lang/elz/src/elz/ast"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBindingRule(t *testing.T) {
	testCases := []struct {
		name            string
		code            string
		expectedBinding *ast.Binding
	}{
		{
			name: "no param bind",
			code: `i = 1`,
			expectedBinding: &ast.Binding{
				Name:      "i",
				ParamList: []string{},
				Expr:      ast.NewInt("1"),
			},
		},
		{
			name: "with param",
			code: `add x y = x + y`,
			expectedBinding: &ast.Binding{
				Name:      "add",
				ParamList: []string{"x", "y"},
				Expr: &ast.BinaryExpr{
					Op:    "+",
					LExpr: ast.NewIdent("x"),
					RExpr: ast.NewIdent("y"),
				},
			},
		},
		{
			name: "function call",
			code: `addOne y = add(1, y)`,
			expectedBinding: &ast.Binding{
				Name:      "addOne",
				ParamList: []string{"y"},
				Expr: &ast.FuncCall{
					AccessChain: "add",
					ArgList: []*ast.Arg{
						{Expr: ast.NewInt("1")},
						{Expr: ast.NewIdent("y")},
					},
				},
			},
		},
		{
			name: "function call with arg name",
			code: `tests = assert(that: 1, should_be: 1)`,
			expectedBinding: &ast.Binding{
				Name:      "tests",
				ParamList: []string{},
				Expr: &ast.FuncCall{
					AccessChain: "assert",
					ArgList: []*ast.Arg{
						{Ident: "that", Expr: ast.NewInt("1")},
						{Ident: "should_be", Expr: ast.NewInt("1")},
					},
				},
			},
		},
		{
			name: "function call as argument",
			code: `tests = add(1, add(add(1, 2), 2))`,
			expectedBinding: &ast.Binding{
				Name:      "tests",
				ParamList: []string{},
				Expr: &ast.FuncCall{
					AccessChain: "add",
					ArgList: []*ast.Arg{
						{Ident: "", Expr: ast.NewInt("1")},
						{Ident: "", Expr: &ast.FuncCall{
							AccessChain: "add",
							ArgList: []*ast.Arg{
								{Ident: "", Expr: &ast.FuncCall{
									AccessChain: "add",
									ArgList: []*ast.Arg{
										{Expr: ast.NewInt("1")},
										{Expr: ast.NewInt("2")},
									},
								}},
								{Expr: ast.NewInt("2")},
							},
						}},
					},
				},
			},
		},
		{
			name: "with certain type limiter",
			code: `
add :: int -> int -> int
add x y = x + y
`,
			expectedBinding: &ast.Binding{
				Name:      "add",
				ParamList: []string{"x", "y"},
				Expr: &ast.BinaryExpr{
					Op:    "+",
					LExpr: ast.NewIdent("x"),
					RExpr: ast.NewIdent("y"),
				},
				Type: []ast.Type{
					&ast.ExistType{Name: "int"},
					&ast.ExistType{Name: "int"},
					&ast.ExistType{Name: "int"},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			builder := New()
			builder.BuildFromCode(testCase.code)
			if testCase.expectedBinding != nil {
				tree := builder.GetTree()
				b, err := tree.GetBinding(testCase.expectedBinding.Name)
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedBinding, b.Binding)
			}
		})
	}
}
