// FILE: internal/lang/object/object.go

package object

import (
	"fmt"
	//"github.com/DauletBai/tenge/internal/lang/ast"
	"github.com/shopspring/decimal"
)

// ObjectType is a string representation of an object's type.
type ObjectType string

// All object types are now based on the tenge language keywords.
const (
	SAN_OBJ    = "SAN"
	AQSHA_OBJ  = "AQSHA"
	AQIQAT_OBJ = "AQIQAT"
	NULL_OBJ   = "NULL"
	QAITAR_VAL = "QAITAR_VAL"
	ERROR_OBJ  = "ERROR"
)

// Singleton instances for common values, named after the language's philosophy.
var (
	NULL = &Null{}
	JAN  = &Aqıqat{Value: true}  // "jan" - soul, represents truth
	JYN  = &Aqıqat{Value: false} // "j'n" - demon, represents falsehood
)

// Object is the interface that every value in tenge will implement.
type Object interface {
	Type() ObjectType
	Inspect() string
}

// --- Object Structs ---

// San represents an integer object.
type San struct {
	Value int64
}
func (s *San) Type() ObjectType { return SAN_OBJ }
func (s *San) Inspect() string  { return fmt.Sprintf("%d", s.Value) }

// Aqsha represents a decimal object for financial calculations.
type Aqsha struct {
	Value decimal.Decimal
}
func (a *Aqsha) Type() ObjectType { return AQSHA_OBJ }
func (a *Aqsha) Inspect() string  { return a.Value.String() }

// Aqıqat represents a boolean object.
type Aqıqat struct {
	Value bool
}
func (a *Aqıqat) Type() ObjectType { return AQIQAT_OBJ }
func (a *Aqıqat) Inspect() string {
	if a.Value {
		return "jan"
	}
	return "j'n"
}

// Null represents the absence of a value.
type Null struct{}
func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// QaıtarValue is a wrapper to handle return values.
type QaıtarValue struct {
	Value Object
}
func (qv *QaıtarValue) Type() ObjectType { return QAITAR_VAL }
func (qv *QaıtarValue) Inspect() string  { return qv.Value.Inspect() }

// Error represents a runtime error.
type Error struct {
	Message string
}
func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "QATE: " + e.Message } // QATE: Kazakh for Error

// --- Environment ---

type Environment struct {
	store map[string]Object
	outer *Environment
}
func NewEnvironment() *Environment { /* ... */ return &Environment{store: make(map[string]Object)} }
func (e *Environment) Get(name string) (Object, bool) { /* ... */ obj, ok := e.store[name]; if !ok && e.outer != nil { obj, ok = e.outer.Get(name) }; return obj, ok }
func (e *Environment) Set(name string, val Object) Object { /* ... */ e.store[name] = val; return val }