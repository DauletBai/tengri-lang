// FILE: internal/lang/evaluator/evaluator.go
// Purpose: Evaluator for Tengri with arrays, calls, and builtin sort().

package evaluator

import (
	"fmt"
	"sort"

	"github.com/DauletBai/tengri-lang/internal/lang/ast"
	"github.com/DauletBai/tengri-lang/internal/lang/object"
)

// Builtins
var builtins = map[string]*object.Builtin{
	"sort": {Fn: builtinSort},
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) { return val }
		return &object.Error{Message: fmt.Sprintf("return %s", val.Inspect())} // simple return signalling

	case *ast.ConstStatement:
		val := Eval(node.Value, env)
		if isError(val) { return val }
		env.Set(node.Name.Value, val)
		return val

	case *ast.JasaStatement:
		val := Eval(node.Value, env)
		if isError(val) { return val }
		env.Set(node.Name.Value, val)
		return val

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		if node.Value { return object.TRUE } ; return object.FALSE

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) { return right }
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) { return left }
		right := Eval(node.Right, env)
		if isError(right) { return right }
		return evalInfixExpression(node.Operator, left, right)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) { return elements[0] }
		return &object.Array{Elements: elements}

	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) { return fn }
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) { return args[0] }
		return applyFunction(fn, args)

	default:
		return object.NULL
	}
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object = object.NULL
	for _, stmt := range program.Statements {
		result = Eval(stmt, env)
		if isError(result) { return result }
	}
	return result
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) { return []object.Object{evaluated} }
		result = append(result, evaluated)
	}
	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if b, ok := builtins[node.Value]; ok {
		return b
	}
	return newError("identifier not found: %s", node.Value)
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return nativeBoolToBooleanObject(!object.IsTruthy(right))
	case "-":
		if r, ok := right.(*object.Integer); ok {
			return &object.Integer{Value: -r.Value}
		}
		return newError("unknown operator: -%s", right.Type())
	default:
		return newError("unknown operator: %s%s", op, right.Type())
	}
}

func evalInfixExpression(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		l := left.(*object.Integer).Value
		r := right.(*object.Integer).Value
		switch op {
		case "+": return &object.Integer{Value: l + r}
		case "-": return &object.Integer{Value: l - r}
		case "*": return &object.Integer{Value: l * r}
		case "/": return &object.Integer{Value: l / r}
		case "==": return nativeBoolToBooleanObject(l == r)
		case "!=": return nativeBoolToBooleanObject(l != r)
		case "<":  return nativeBoolToBooleanObject(l < r)
		case "<=": return nativeBoolToBooleanObject(l <= r)
		case ">":  return nativeBoolToBooleanObject(l > r)
		case ">=": return nativeBoolToBooleanObject(l >= r)
		}
	case op == "==":
		return nativeBoolToBooleanObject(left == right)
	case op == "!=":
		return nativeBoolToBooleanObject(left != right)
	}
	return newError("type mismatch: %s %s %s", left.Type(), op, right.Type())
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch f := fn.(type) {
	case *object.Builtin:
		return f.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func builtinSort(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("sort expects exactly 1 argument (array)")
	}
	arr, ok := args[0].(*object.Array)
	if !ok {
		return newError("sort expects array, got %s", args[0].Type())
	}
	// Copy and sort; only integers supported for benchmark simplicity
	copyElems := make([]object.Object, len(arr.Elements))
	copy(copyElems, arr.Elements)
	ints := make([]int, len(copyElems))
	for i, el := range copyElems {
		iv, ok := el.(*object.Integer)
		if !ok { return newError("sort supports only integer arrays") }
		ints[i] = int(iv.Value)
	}
	sort.Ints(ints)
	res := make([]object.Object, len(ints))
	for i, v := range ints { res[i] = &object.Integer{Value: int64(v)} }
	return &object.Array{Elements: res}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input { return object.TRUE } ; return object.FALSE
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}
	return false
}