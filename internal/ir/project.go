package ir

import "github.com/shishir1290/goduck/internal/ast"

type Project struct {
	AppName string
	Classes []*ast.ClassDeclaration
}