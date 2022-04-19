package luadump

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gueckmooh/bs/pkg/functional"
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/ast"
)

func DumpFunction(F *lua.LFunction) string {
	// return DumpFunctionExpr(F.Proto.Chunk, 0)
	d := NewDumper(0)
	return d.DumpExpr(F.Proto.Chunk)
}

func NewDumper(n int) *Dumper {
	v := new(Dumper)
	*v = Dumper(n)
	return v
}

func (d *Dumper) indent() string {
	return strings.Repeat("    ", int(*d))
}

func (d *Dumper) inc() {
	*d = *d + 1
}

func (d *Dumper) deinc() {
	*d = *d - 1
}

type Dumper int

func (d *Dumper) DumpStmt(s ast.Stmt) string {
	switch stmt := s.(type) {
	case *ast.FuncCallStmt:
		return fmt.Sprintf("%s%s", d.indent(), d.DumpExpr(stmt.Expr))
	case *ast.AssignStmt:
		return fmt.Sprintf("%s%s = %s", d.indent(),
			strings.Join(functional.ListMap(stmt.Lhs, d.DumpExpr), ", "),
			strings.Join(functional.ListMap(stmt.Rhs, d.DumpExpr), ", "))
	default:
		return fmt.Sprintf("stmt: %#v", stmt)
	}
}

func (d *Dumper) DumpExpr(e ast.Expr) string {
	switch expr := e.(type) {
	case *ast.FuncCallExpr:
		return fmt.Sprintf("%s(%s)", d.DumpExpr(expr.Func),
			strings.Join(functional.ListMap(expr.Args, d.DumpExpr), ", "))
	case *ast.FunctionExpr:
		{
			var buf bytes.Buffer
			fmt.Fprintf(&buf, "function (%s)\n", strings.Join(expr.ParList.Names, ", "))
			d.inc()
			for _, stmt := range expr.Stmts {
				fmt.Fprintf(&buf, "%s\n", d.DumpStmt(stmt))
			}
			d.deinc()
			fmt.Fprintf(&buf, "end")
			return buf.String()
		}
	case *ast.IdentExpr:
		return fmt.Sprintf("%s", expr.Value)
	case *ast.StringExpr:
		return fmt.Sprintf(`"%s"`, expr.Value)
	case *ast.AttrGetExpr:
		return fmt.Sprintf("%s[%s]", d.DumpExpr(expr.Object), d.DumpExpr(expr.Key))
	case *ast.StringConcatOpExpr:
		return fmt.Sprintf("%s .. %s", d.DumpExpr(expr.Lhs), d.DumpExpr(expr.Rhs))
	default:
		return fmt.Sprintf("expr: %#v", expr)
	}
}
