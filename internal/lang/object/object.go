// FILE: internal/lang/object/object.go
// Purpose: Runtime object system for Tengri with arrays and builtins; Kazakh boolean Inspect.

package object

import (
	"bytes"
	"fmt"
	//"strings"

	//"github.com/DauletBai/tengri-lang/internal/lang/ast"
)

type ObjectType string

const (
	INTEGER = "INTEGER"
	BOOLEAN = "BOOLEAN"
	ARRAY   = "ARRAY"
	NULL_OBJ    = "NULL"
	ERROR   = "ERROR"
	BUILTIN = "BUILTIN"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// --- Primitives ---

type Integer struct { 
	Value int64 
}

func (i *Integer) Type() ObjectType   { return INTEGER }
func (i *Integer) Inspect() string    { return fmt.Sprintf("%d", i.Value) }

type Boolean struct { 
	Value bool 
}

func (b *Boolean) Type() ObjectType   { return BOOLEAN }
func (b *Boolean) Inspect() string    { if b.Value { return "jan" } ; return "j'n" }

type Null struct{}
func (n *Null) Type() ObjectType      { return NULL_OBJ }
func (n *Null) Inspect() string       { return "" }

// --- Composite ---

type Array struct { 
	Elements []Object 
}

func (a *Array) Type() ObjectType     { return ARRAY }
func (a *Array) Inspect() string {
	var out bytes.Buffer
	out.WriteString("[")
	for i, e := range a.Elements {
		if i > 0 { out.WriteString(", ") }
		out.WriteString(e.Inspect())
	}
	out.WriteString("]")
	return out.String()
}

// --- Error ---
type Error struct { 
	Message string 
}

func (e *Error) Type() ObjectType     { return ERROR }
func (e *Error) Inspect() string      { return "ERROR: " + e.Message }

// --- Builtin ---
type BuiltinFunction func(args ...Object) Object
type Builtin struct { 
	Fn BuiltinFunction 
}

func (b *Builtin) Type() ObjectType   { return BUILTIN }
func (b *Builtin) Inspect() string    { return "<builtin>" }

// --- Environment ---
type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object)}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// Singletons
var (
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
	NULL  = &Null{}
)

func IsTruthy(obj Object) bool {
	switch o := obj.(type) {
	case *Null:
		return false
	case *Boolean:
		return o.Value
	default:
		return true
	}
}