package semantic

import (
	"fmt"
	"strings"

	"github.com/shishir1290/goduck/internal/ast"
)

type Analyzer struct {
	program *ast.Program
	classes map[string]*ClassInfo
}

type ClassInfo struct {
	Name    string
	Fields  map[string]string // name -> type name
	Methods map[string]*MethodInfo
}

type MethodInfo struct {
	Name       string
	Parameters []ParamInfo
	ReturnType string
}

type ParamInfo struct {
	Name string
	Type string
}

type Env struct {
	outer *Env
	vars  map[string]string
}

func NewEnv(outer *Env) *Env {
	return &Env{
		outer: outer,
		vars:  make(map[string]string),
	}
}

func (e *Env) Set(name, typeStr string) {
	e.vars[name] = typeStr
}

func (e *Env) Get(name string) (string, bool) {
	if t, ok := e.vars[name]; ok {
		return t, true
	}
	if e.outer != nil {
		return e.outer.Get(name)
	}
	return "", false
}

func New(program *ast.Program) *Analyzer {
	return &Analyzer{
		program: program,
		classes: make(map[string]*ClassInfo),
	}
}

func (a *Analyzer) Analyze() error {
	// Step 1: Validate App
	if a.program.App == nil {
		return fmt.Errorf("missing app declaration")
	}
	if a.program.App.Name == "" {
		return fmt.Errorf("app name cannot be empty")
	}

	// Step 2: Register all classes
	for _, class := range a.program.Classes {
		if _, exists := a.classes[class.Name]; exists {
			return fmt.Errorf("duplicate class declaration: %s", class.Name)
		}

		info := &ClassInfo{
			Name:    class.Name,
			Fields:  make(map[string]string),
			Methods: make(map[string]*MethodInfo),
		}

		for _, field := range class.Fields {
			if field.Type == nil {
				info.Fields[field.Name] = "any"
			} else {
				typeName := field.Type.Name
				if field.Type.IsArray {
					typeName += "[]"
				}
				info.Fields[field.Name] = typeName
			}
		}

		for _, method := range class.Methods {
			mInfo := &MethodInfo{
				Name:       method.Name,
				Parameters: []ParamInfo{},
				ReturnType: "void",
			}

			if method.ReturnType != nil {
				retName := method.ReturnType.Name
				if method.ReturnType.IsArray {
					retName += "[]"
				}
				mInfo.ReturnType = retName
			}

			for _, param := range method.Parameters {
				pName := param.Type.Name
				if param.Type.IsArray {
					pName += "[]"
				}
				mInfo.Parameters = append(mInfo.Parameters, ParamInfo{
					Name: param.Name,
					Type: pName,
				})
			}

			info.Methods[method.Name] = mInfo
		}

		a.classes[class.Name] = info
	}

	// Step 3: Type check global variables
	globalEnv := NewEnv(nil)
	for _, v := range a.program.Variables {
		if v.Value != nil {
			valType, err := a.evaluateType(v.Value, globalEnv, "")
			if err != nil {
				return err
			}
			if v.Type != nil {
				expectedType := v.Type.Name
				if v.Type.IsArray {
					expectedType += "[]"
				}
				if !a.typesCompatible(expectedType, valType) {
					return fmt.Errorf("variable '%s' type mismatch: declared as %s, assigned %s", v.Name, expectedType, valType)
				}
				globalEnv.Set(v.Name, expectedType)
			} else {
				globalEnv.Set(v.Name, valType)
			}
		} else if v.Type != nil {
			expectedType := v.Type.Name
			if v.Type.IsArray {
				expectedType += "[]"
			}
			globalEnv.Set(v.Name, expectedType)
		}
	}

	// Step 4: Type check all class methods
	for _, class := range a.program.Classes {
		for _, method := range class.Methods {
			methodEnv := NewEnv(globalEnv)
			info := a.classes[class.Name]
			mInfo := info.Methods[method.Name]

			// Add method parameters to scope
			for _, param := range mInfo.Parameters {
				methodEnv.Set(param.Name, param.Type)
			}

			// Check method body
			if method.Body != nil {
				if err := a.checkBlock(method.Body, methodEnv, class.Name, mInfo.ReturnType); err != nil {
					return fmt.Errorf("class %s method %s: %w", class.Name, method.Name, err)
				}
			}
		}
	}

	return nil
}

func (a *Analyzer) checkBlock(block *ast.BlockStatement, env *Env, currentClass string, expectedReturn string) error {
	for _, stmt := range block.Statements {
		if err := a.checkStatement(stmt, env, currentClass, expectedReturn); err != nil {
			return err
		}
	}
	return nil
}

func (a *Analyzer) checkStatement(stmt ast.Statement, env *Env, currentClass string, expectedReturn string) error {
	switch s := stmt.(type) {
	case *ast.VariableDeclaration:
		if s.Value != nil {
			valType, err := a.evaluateType(s.Value, env, currentClass)
			if err != nil {
				return err
			}
			if s.Type != nil {
				expectedType := s.Type.Name
				if s.Type.IsArray {
					expectedType += "[]"
				}
				if !a.typesCompatible(expectedType, valType) {
					return fmt.Errorf("variable '%s' type mismatch: declared as %s, assigned %s", s.Name, expectedType, valType)
				}
				env.Set(s.Name, expectedType)
			} else {
				env.Set(s.Name, valType)
			}
		} else if s.Type != nil {
			expectedType := s.Type.Name
			if s.Type.IsArray {
				expectedType += "[]"
			}
			env.Set(s.Name, expectedType)
		}

	case *ast.ReturnStatement:
		valType := "void"
		if s.Value != nil {
			var err error
			valType, err = a.evaluateType(s.Value, env, currentClass)
			if err != nil {
				return err
			}
		}
		if !a.typesCompatible(expectedReturn, valType) {
			return fmt.Errorf("return type mismatch: method expects %s, returned %s", expectedReturn, valType)
		}

	case *ast.ExpressionStatement:
		_, err := a.evaluateType(s.Expression, env, currentClass)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Analyzer) evaluateType(expr ast.Expression, env *Env, currentClass string) (string, error) {
	if expr == nil {
		return "void", nil
	}

	switch e := expr.(type) {
	case *ast.Literal:
		return e.Type, nil

	case *ast.Identifier:
		if e.Name == "true" || e.Name == "false" {
			return "boolean", nil
		}
		if e.Name == "this" {
			if currentClass == "" {
				return "", fmt.Errorf("use of 'this' outside of a class method")
			}
			return currentClass, nil
		}
		if t, ok := env.Get(e.Name); ok {
			return t, nil
		}
		return "", fmt.Errorf("undefined identifier: %s", e.Name)

	case *ast.MemberExpression:
		objType, err := a.evaluateType(e.Object, env, currentClass)
		if err != nil {
			return "", err
		}

		classInfo, ok := a.classes[objType]
		if !ok {
			return "", fmt.Errorf("type '%s' has no members", objType)
		}

		// Check field
		if fieldType, ok := classInfo.Fields[e.Property]; ok {
			return fieldType, nil
		}

		// Check method
		if methodInfo, ok := classInfo.Methods[e.Property]; ok {
			return "method:" + methodInfo.ReturnType, nil
		}

		return "", fmt.Errorf("member '%s' not found on class '%s'", e.Property, objType)

	case *ast.CallExpression:
		var returnType string
		var paramTypes []string

		switch fn := e.Function.(type) {
		case *ast.MemberExpression:
			objType, err := a.evaluateType(fn.Object, env, currentClass)
			if err != nil {
				return "", err
			}

			classInfo, ok := a.classes[objType]
			if !ok {
				return "", fmt.Errorf("type '%s' has no methods", objType)
			}

			methodInfo, ok := classInfo.Methods[fn.Property]
			if !ok {
				return "", fmt.Errorf("method '%s' not found on class '%s'", fn.Property, objType)
			}

			returnType = methodInfo.ReturnType
			for _, p := range methodInfo.Parameters {
				paramTypes = append(paramTypes, p.Type)
			}

		case *ast.Identifier:
			if currentClass != "" {
				classInfo := a.classes[currentClass]
				if methodInfo, ok := classInfo.Methods[fn.Name]; ok {
					returnType = methodInfo.ReturnType
					for _, p := range methodInfo.Parameters {
						paramTypes = append(paramTypes, p.Type)
					}
				} else {
					return "", fmt.Errorf("undefined method: %s", fn.Name)
				}
			} else {
				return "", fmt.Errorf("undefined function: %s", fn.Name)
			}

		default:
			return "", fmt.Errorf("invalid call expression target")
		}

		if len(e.Arguments) != len(paramTypes) {
			return "", fmt.Errorf("expected %d arguments, got %d", len(paramTypes), len(e.Arguments))
		}

		for i, arg := range e.Arguments {
			argType, err := a.evaluateType(arg, env, currentClass)
			if err != nil {
				return "", err
			}
			expectedType := paramTypes[i]
			if !a.typesCompatible(expectedType, argType) {
				return "", fmt.Errorf("argument %d type mismatch: expected %s, got %s", i+1, expectedType, argType)
			}
		}

		return returnType, nil

	case *ast.BinaryExpression:
		leftType, err := a.evaluateType(e.Left, env, currentClass)
		if err != nil {
			return "", err
		}
		rightType, err := a.evaluateType(e.Right, env, currentClass)
		if err != nil {
			return "", err
		}

		if e.Operator == "+" {
			if leftType == "string" || rightType == "string" {
				return "string", nil
			}
			if leftType == "number" && rightType == "number" {
				return "number", nil
			}
			return "", fmt.Errorf("operator '+' cannot be applied to types '%s' and '%s'", leftType, rightType)
		}

		return "", fmt.Errorf("unknown operator: %s", e.Operator)
	}

	return "any", nil
}

func (a *Analyzer) typesCompatible(expected, actual string) bool {
	if expected == "any" || actual == "any" {
		return true
	}
	if expected == actual {
		return true
	}
	// Promise compatibility
	if strings.HasPrefix(expected, "Promise<") && strings.HasSuffix(expected, ">") {
		inner := expected[len("Promise<") : len(expected)-1]
		if inner == actual {
			return true
		}
	}
	if strings.HasPrefix(actual, "Promise<") && strings.HasSuffix(actual, ">") {
		inner := actual[len("Promise<") : len(actual)-1]
		if inner == expected {
			return true
		}
	}
	return false
}