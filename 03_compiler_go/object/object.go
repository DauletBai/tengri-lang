// 03_compiler_go/object/object.go
package object

import (
	"bytes"
	"fmt"
	"strings"
	"tengri-lang/03_compiler_go/ast"
)

type ObjectType string

const (
	INTEGER  = "INTEGER"
	BOOLEAN  = "BOOLEAN"
	NULL_OBJ = "NULL"
	FUNCTION = "FUNCTION"
	RETURN   = "RETURN"
	ERROR    = "ERROR"
)

var (
	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}
func (i *Integer) Type() ObjectType { return INTEGER }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}
func (b *Boolean) Type() ObjectType { return BOOLEAN }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct{}
func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type Function struct {
	Parameters []*ast.Parameter
	Body       *ast.BlockStatement
	Env        *Environment
}
func (f *Function) Type() ObjectType { return FUNCTION }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {...}")
	return out.String()
}

type ReturnValue struct {
	Value Object
}
func (rv *ReturnValue) Type() ObjectType { return RETURN }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}
func (e *Error) Type() ObjectType { return ERROR }
func (e *Error) Inspect() string  { return "ОШИБКА: " + e.Message }