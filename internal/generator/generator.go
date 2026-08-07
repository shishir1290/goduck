package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shishir1290/goduck/internal/ast"
)

type Generator struct {
	program *ast.Program
	project *Project
}

func New(program *ast.Program) *Generator {
	name := "App"
	if program.App != nil {
		name = program.App.Name
	}
	return &Generator{
		program: program,
		project: &Project{
			Name: name,
		},
	}
}

func (g *Generator) Generate() *Project {
	g.generateGoMod()
	g.generateModels()
	g.generateServices()
	g.generateControllers()
	g.generateModules()
	g.generateMain()

	return g.project
}

func (g *Generator) generateGoMod() {
	relPath := getRelativePathToGoduck(g.project.Name)
	content := fmt.Sprintf(`module %s

go 1.20

require github.com/shishir1290/goduck v0.0.0

replace github.com/shishir1290/goduck => %s
`, strings.ToLower(g.project.Name), relPath)

	g.project.AddFile("go.mod", content)
}

func findGoduckRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if content, err := os.ReadFile(goModPath); err == nil {
			if strings.Contains(string(content), "module github.com/shishir1290/goduck") {
				return dir
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func getRelativePathToGoduck(appName string) string {
	goduckRoot := findGoduckRoot()
	if goduckRoot == "" {
		return "../.."
	}
	absGoduckRoot, err := filepath.Abs(goduckRoot)
	if err != nil {
		return "../.."
	}
	buildDir := filepath.Join("build", appName)
	absBuildDir, err := filepath.Abs(buildDir)
	if err != nil {
		return "../.."
	}
	rel, err := filepath.Rel(absBuildDir, absGoduckRoot)
	if err != nil {
		return filepath.ToSlash(absGoduckRoot)
	}
	return filepath.ToSlash(rel)
}

func (g *Generator) generateModels() {
	var sb strings.Builder
	sb.WriteString("package main\n\n")

	hasModels := false
	for _, class := range g.program.Classes {
		if !hasDecorator(class.Decorators, "@Controller") && !hasDecorator(class.Decorators, "@Service") && !hasDecorator(class.Decorators, "@Module") {
			hasModels = true
			sb.WriteString(fmt.Sprintf("type %s struct {\n", capitalize(class.Name)))
			for _, field := range class.Fields {
				goType := g.toGoType(field.Type)
				sb.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\"`\n", capitalize(field.Name), goType, field.Name))
			}
			sb.WriteString("}\n\n")
		}
	}

	if hasModels {
		g.project.AddFile("models.go", sb.String())
	}
}

func (g *Generator) generateServices() {
	var sb strings.Builder
	sb.WriteString("package main\n\n")

	hasServices := false
	for _, class := range g.program.Classes {
		if hasDecorator(class.Decorators, "@Service") {
			hasServices = true
			sb.WriteString(fmt.Sprintf("type %s struct {\n", capitalize(class.Name)))
			// Service fields
			for _, field := range class.Fields {
				goType := g.toGoType(field.Type)
				sb.WriteString(fmt.Sprintf("\t%s %s\n", capitalize(field.Name), goType))
			}
			sb.WriteString("}\n\n")

			// Methods
			for _, method := range class.Methods {
				var params []string
				for _, param := range method.Parameters {
					params = append(params, param.Name+" "+g.toGoType(param.Type))
				}
				retType := g.toGoType(method.ReturnType)
				retSig := ""
				if retType != "" {
					retSig = " " + retType
				}

				sb.WriteString(fmt.Sprintf("func (s *%s) %s(%s)%s {\n",
					capitalize(class.Name),
					capitalize(method.Name),
					strings.Join(params, ", "),
					retSig,
				))

				if method.Body != nil {
					sb.WriteString(g.translateBlock(method.Body, "s"))
				} else {
					if retType != "" {
						sb.WriteString("\treturn nil\n")
					}
				}
				sb.WriteString("}\n\n")
			}
		}
	}

	if hasServices {
		g.project.AddFile("services.go", sb.String())
	}
}

func (g *Generator) generateControllers() {
	// Collect services to detect DI fields
	services := make(map[string]bool)
	for _, class := range g.program.Classes {
		if hasDecorator(class.Decorators, "@Service") {
			services[class.Name] = true
		}
	}

	needStrconv := false
	for _, class := range g.program.Classes {
		if hasDecorator(class.Decorators, "@Controller") {
			dec := getDecorator(class.Decorators, "@Controller")
			prefix := ""
			if len(dec.Arguments) > 0 {
				if lit, ok := dec.Arguments[0].(*ast.Literal); ok && lit.Type == "string" {
					prefix = lit.Value
				}
			}

			for _, method := range class.Methods {
				routeDec := getRouteDecorator(method.Decorators)
				if routeDec != nil {
					pathSuffix := ""
					if len(routeDec.Arguments) > 0 {
						if lit, ok := routeDec.Arguments[0].(*ast.Literal); ok && lit.Type == "string" {
							pathSuffix = lit.Value
						}
					}

					fullPath := joinPaths(prefix, pathSuffix)
					for _, param := range method.Parameters {
						if strings.Contains(fullPath, ":"+param.Name) && param.Type.Name == "number" {
							needStrconv = true
						}
					}
				}
			}
		}
	}

	var sb strings.Builder
	sb.WriteString("package main\n\n")
	sb.WriteString("import (\n")
	if needStrconv {
		sb.WriteString("\t\"strconv\"\n")
	}
	sb.WriteString("\t\"github.com/shishir1290/goduck/runtime/router\"\n")
	sb.WriteString("\t\"github.com/shishir1290/goduck/runtime/context\"\n")
	sb.WriteString(")\n\n")

	hasControllers := false
	for _, class := range g.program.Classes {
		if hasDecorator(class.Decorators, "@Controller") {
			hasControllers = true
			sb.WriteString(fmt.Sprintf("type %s struct {\n", capitalize(class.Name)))

			for _, field := range class.Fields {
				goType := g.toGoType(field.Type)
				sb.WriteString(fmt.Sprintf("\t%s %s\n", capitalize(field.Name), goType))
			}
			sb.WriteString("}\n\n")

			// Controller Actions (without routing decorators or standard internal ones)
			for _, method := range class.Methods {
				// Even if they have routing decorator, generate the Go method on the struct
				var params []string
				for _, param := range method.Parameters {
					params = append(params, param.Name+" "+g.toGoType(param.Type))
				}
				retType := g.toGoType(method.ReturnType)
				retSig := ""
				if retType != "" {
					retSig = " " + retType
				}

				sb.WriteString(fmt.Sprintf("func (c *%s) %s(%s)%s {\n",
					capitalize(class.Name),
					capitalize(method.Name),
					strings.Join(params, ", "),
					retSig,
				))

				if method.Body != nil {
					sb.WriteString(g.translateBlock(method.Body, "c"))
				} else {
					if retType != "" {
						sb.WriteString("\treturn nil\n")
					}
				}
				sb.WriteString("}\n\n")
			}

			// Generate SetupRoutes method for the controller
			sb.WriteString(fmt.Sprintf("func (c *%s) SetupRoutes(g *router.Group) {\n", capitalize(class.Name)))
			dec := getDecorator(class.Decorators, "@Controller")
			prefix := ""
			if len(dec.Arguments) > 0 {
				if lit, ok := dec.Arguments[0].(*ast.Literal); ok && lit.Type == "string" {
					prefix = lit.Value
				}
			}

			for _, method := range class.Methods {
				routeDec := getRouteDecorator(method.Decorators)
				if routeDec != nil {
					httpMethod := strings.TrimPrefix(routeDec.Name, "@") // GET, POST, etc.
					pathSuffix := ""
					if len(routeDec.Arguments) > 0 {
						if lit, ok := routeDec.Arguments[0].(*ast.Literal); ok && lit.Type == "string" {
							pathSuffix = lit.Value
						}
					}

					fullPath := joinPaths(prefix, pathSuffix)

					sb.WriteString(fmt.Sprintf("\tg.%s(\"%s\", func(ctx *context.Context) {\n", httpMethod, fullPath))

					// Generate handler body
					var handlerBody []string
					var callArgs []string

					for _, param := range method.Parameters {
						paramGoType := g.toGoType(param.Type)
						if strings.Contains(fullPath, ":"+param.Name) {
							// Path parameter
							if param.Type.Name == "number" {
								handlerBody = append(handlerBody, fmt.Sprintf("\t\t%sStr := ctx.Param(\"%s\")", param.Name, param.Name))
								handlerBody = append(handlerBody, fmt.Sprintf("\t\t%s, _ := strconv.Atoi(%sStr)", param.Name, param.Name))
								callArgs = append(callArgs, param.Name)
							} else {
								handlerBody = append(handlerBody, fmt.Sprintf("\t\t%s := ctx.Param(\"%s\")", param.Name, param.Name))
								callArgs = append(callArgs, param.Name)
							}
						} else {
							// Request Body parameter
							baseType := strings.TrimPrefix(paramGoType, "*")
							handlerBody = append(handlerBody, fmt.Sprintf("\t\tvar %s %s", param.Name, baseType))
							handlerBody = append(handlerBody, fmt.Sprintf("\t\tif err := ctx.Bind(&%s); err != nil {", param.Name))
							handlerBody = append(handlerBody, "\t\t\tctx.Status(400)")
							handlerBody = append(handlerBody, "\t\t\tctx.String(err.Error())")
							handlerBody = append(handlerBody, "\t\t\treturn")
							handlerBody = append(handlerBody, "\t\t}")
							if strings.HasPrefix(paramGoType, "*") {
								callArgs = append(callArgs, "&"+param.Name)
							} else {
								callArgs = append(callArgs, param.Name)
							}
						}
					}

					methodCall := fmt.Sprintf("c.%s(%s)", capitalize(method.Name), strings.Join(callArgs, ", "))
					retType := g.toGoType(method.ReturnType)

					for _, line := range handlerBody {
						sb.WriteString(line + "\n")
					}

					if retType == "" {
						sb.WriteString(fmt.Sprintf("\t\t%s\n", methodCall))
						sb.WriteString("\t\tctx.Status(200)\n")
					} else if retType == "string" {
						sb.WriteString(fmt.Sprintf("\t\tres := %s\n", methodCall))
						sb.WriteString("\t\tctx.String(res)\n")
					} else {
						sb.WriteString(fmt.Sprintf("\t\tres := %s\n", methodCall))
						sb.WriteString("\t\tctx.JSON(res)\n")
					}

					sb.WriteString("\t})\n")
				}
			}
			sb.WriteString("}\n\n")
		}
	}

	if hasControllers {
		g.project.AddFile("controllers.go", sb.String())
	}
}

func (g *Generator) generateMain() {
	// Collect services to detect DI fields
	services := make(map[string]bool)
	for _, class := range g.program.Classes {
		if hasDecorator(class.Decorators, "@Service") {
			services[class.Name] = true
		}
	}

	var sb strings.Builder
	sb.WriteString("package main\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"log\"\n")
	sb.WriteString("\n\t\"github.com/shishir1290/goduck/runtime/server\"\n")
	sb.WriteString("\t\"github.com/shishir1290/goduck/runtime/container\"\n")
	sb.WriteString(")\n\n")

	sb.WriteString("func main() {\n")

	// Get port from App config
	port := 8080
	if g.program.App != nil {
		for _, prop := range g.program.App.Properties {
			if prop.Name == "port" {
				if lit, ok := prop.Value.(*ast.Literal); ok && lit.Type == "number" {
					fmt.Sscanf(lit.Value, "%d", &port)
				}
			}
		}
	}

	sb.WriteString(fmt.Sprintf("\ts := server.New(%d)\n", port))
	sb.WriteString("\tctr := s.Container()\n\n")

	// Register services in DI container
	for _, class := range g.program.Classes {
		if hasDecorator(class.Decorators, "@Service") {
			var resolvedFields []string
			for _, field := range class.Fields {
				goType := g.toGoType(field.Type)
				typeName := strings.TrimPrefix(goType, "*")
				if services[typeName] {
					resolvedFields = append(resolvedFields, fmt.Sprintf("%s: c.Resolve((*%s)(nil)).(*%s)", capitalize(field.Name), typeName, typeName))
				}
			}

			sb.WriteString(fmt.Sprintf("\tctr.Register((*%s)(nil), func(c *container.Container) any {\n", capitalize(class.Name)))
			if len(resolvedFields) > 0 {
				sb.WriteString(fmt.Sprintf("\t\treturn &%s{\n", capitalize(class.Name)))
				for _, rf := range resolvedFields {
					sb.WriteString(fmt.Sprintf("\t\t\t%s,\n", rf))
				}
				sb.WriteString("\t\t}\n")
			} else {
				sb.WriteString(fmt.Sprintf("\t\treturn &%s{}\n", capitalize(class.Name)))
			}
			sb.WriteString("\t})\n")
		}
	}

	// Register controllers in DI container
	for _, class := range g.program.Classes {
		if hasDecorator(class.Decorators, "@Controller") {
			var resolvedFields []string
			for _, field := range class.Fields {
				goType := g.toGoType(field.Type)
				typeName := strings.TrimPrefix(goType, "*")
				if services[typeName] {
					resolvedFields = append(resolvedFields, fmt.Sprintf("%s: c.Resolve((*%s)(nil)).(*%s)", capitalize(field.Name), typeName, typeName))
				}
			}

			sb.WriteString(fmt.Sprintf("\tctr.Register((*%s)(nil), func(c *container.Container) any {\n", capitalize(class.Name)))
			if len(resolvedFields) > 0 {
				sb.WriteString(fmt.Sprintf("\t\treturn &%s{\n", capitalize(class.Name)))
				for _, rf := range resolvedFields {
					sb.WriteString(fmt.Sprintf("\t\t\t%s,\n", rf))
				}
				sb.WriteString("\t\t}\n")
			} else {
				sb.WriteString(fmt.Sprintf("\t\treturn &%s{}\n", capitalize(class.Name)))
			}
			sb.WriteString("\t})\n")
		}
	}

	// Register modules in DI container
	for _, class := range g.program.Classes {
		if hasDecorator(class.Decorators, "@Module") {
			var resolvedFields []string
			for _, field := range class.Fields {
				goType := g.toGoType(field.Type)
				typeName := strings.TrimPrefix(goType, "*")
				resolvedFields = append(resolvedFields, fmt.Sprintf("%s: c.Resolve((*%s)(nil)).(*%s)", capitalize(field.Name), typeName, typeName))
			}

			sb.WriteString(fmt.Sprintf("\tctr.Register((*%s)(nil), func(c *container.Container) any {\n", capitalize(class.Name)))
			if len(resolvedFields) > 0 {
				sb.WriteString(fmt.Sprintf("\t\treturn &%s{\n", capitalize(class.Name)))
				for _, rf := range resolvedFields {
					sb.WriteString(fmt.Sprintf("\t\t\t%s,\n", rf))
				}
				sb.WriteString("\t\t}\n")
			} else {
				sb.WriteString(fmt.Sprintf("\t\treturn &%s{}\n", capitalize(class.Name)))
			}
			sb.WriteString("\t})\n")
		}
	}
	sb.WriteString("\n")

	// Set up routing via Modules
	apiPath := "/api"
	if g.program.App != nil {
		for _, prop := range g.program.App.Properties {
			if prop.Name == "apiPath" {
				if lit, ok := prop.Value.(*ast.Literal); ok && lit.Type == "string" {
					apiPath = lit.Value
				}
			}
		}
	}
	sb.WriteString(fmt.Sprintf("\tapiGroup := s.Group(\"%s\")\n\n", apiPath))

	for _, class := range g.program.Classes {
		if hasDecorator(class.Decorators, "@Module") {
			moduleVar := strings.ToLower(class.Name)
			sb.WriteString(fmt.Sprintf("\t%s := ctr.Resolve((*%s)(nil)).(*%s)\n", moduleVar, capitalize(class.Name), capitalize(class.Name)))
			sb.WriteString(fmt.Sprintf("\t%s.SetupRoutes(apiGroup)\n", moduleVar))
		}
	}
	sb.WriteString("\n")

	sb.WriteString("\tlog.Fatal(s.Run())\n")
	sb.WriteString("}\n\n")

	g.project.AddFile("main.go", sb.String())
}


// Helpers
func hasDecorator(decs []*ast.Decorator, name string) bool {
	for _, dec := range decs {
		if dec.Name == name {
			return true
		}
	}
	return false
}

func getDecorator(decs []*ast.Decorator, name string) *ast.Decorator {
	for _, dec := range decs {
		if dec.Name == name {
			return dec
		}
	}
	return nil
}

func getRouteDecorator(decs []*ast.Decorator) *ast.Decorator {
	routes := map[string]bool{
		"@GET":    true,
		"@POST":   true,
		"@PUT":    true,
		"@PATCH":  true,
		"@DELETE": true,
	}
	for _, dec := range decs {
		if routes[dec.Name] {
			return dec
		}
	}
	return nil
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func joinPaths(prefix, suffix string) string {
	p := "/" + strings.Trim(prefix, "/")
	s := strings.Trim(suffix, "/")
	if s == "" {
		if p == "/" {
			return "/"
		}
		return p
	}
	if p == "/" {
		return "/" + s
	}
	return p + "/" + s
}

func (g *Generator) toGoType(t *ast.TypeNode) string {
	if t == nil {
		return ""
	}
	name := t.Name
	if strings.HasPrefix(name, "Promise<") && strings.HasSuffix(name, ">") {
		name = name[len("Promise<") : len(name)-1]
	}

	var goType string
	switch name {
	case "string":
		goType = "string"
	case "number":
		goType = "int"
	case "boolean":
		goType = "bool"
	case "void":
		goType = ""
	default:
		goType = "*" + capitalize(name)
	}

	if t.IsArray {
		goType = "[]" + goType
	}
	return goType
}

func (g *Generator) translateBlock(block *ast.BlockStatement, receiver string) string {
	var lines []string
	for _, stmt := range block.Statements {
		lines = append(lines, "\t"+g.translateStatement(stmt, receiver))
	}
	return strings.Join(lines, "\n") + "\n"
}

func (g *Generator) translateStatement(stmt ast.Statement, receiver string) string {
	switch s := stmt.(type) {
	case *ast.ReturnStatement:
		if s.Value == nil {
			return "return"
		}
		return "return " + g.translateExpr(s.Value, receiver)
	case *ast.VariableDeclaration:
		if s.Value != nil {
			return s.Name + " := " + g.translateExpr(s.Value, receiver)
		}
		goType := g.toGoType(s.Type)
		return "var " + s.Name + " " + goType
	case *ast.ExpressionStatement:
		return g.translateExpr(s.Expression, receiver)
	}
	return ""
}

func (g *Generator) translateExpr(expr ast.Expression, receiver string) string {
	if expr == nil {
		return ""
	}

	switch e := expr.(type) {
	case *ast.Literal:
		if e.Type == "string" {
			return "\"" + e.Value + "\""
		}
		return e.Value

	case *ast.Identifier:
		if e.Name == "this" {
			return receiver
		}
		return e.Name

	case *ast.MemberExpression:
		obj := g.translateExpr(e.Object, receiver)
		prop := capitalize(e.Property)
		return obj + "." + prop

	case *ast.CallExpression:
		fn := g.translateExpr(e.Function, receiver)
		var args []string
		for _, arg := range e.Arguments {
			args = append(args, g.translateExpr(arg, receiver))
		}
		return fn + "(" + strings.Join(args, ", ") + ")"

	case *ast.BinaryExpression:
		left := g.translateExpr(e.Left, receiver)
		right := g.translateExpr(e.Right, receiver)
		return left + " " + e.Operator + " " + right
	}

	return ""
}

func (g *Generator) generateModules() {
	var sb strings.Builder
	sb.WriteString("package main\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"github.com/shishir1290/goduck/runtime/router\"\n")
	sb.WriteString(")\n\n")

	hasModules := false
	for _, class := range g.program.Classes {
		if hasDecorator(class.Decorators, "@Module") {
			hasModules = true
			sb.WriteString(fmt.Sprintf("type %s struct {\n", capitalize(class.Name)))
			for _, field := range class.Fields {
				goType := g.toGoType(field.Type)
				sb.WriteString(fmt.Sprintf("\t%s %s\n", capitalize(field.Name), goType))
			}
			sb.WriteString("}\n\n")

			sb.WriteString(fmt.Sprintf("func (m *%s) SetupRoutes(g *router.Group) {\n", capitalize(class.Name)))
			// Find controller fields in this module to call SetupRoutes on them
			for _, field := range class.Fields {
				goType := g.toGoType(field.Type)
				typeName := strings.TrimPrefix(goType, "*")
				// Check if the referenced class has a @Controller decorator
				for _, otherClass := range g.program.Classes {
					if otherClass.Name == typeName && hasDecorator(otherClass.Decorators, "@Controller") {
						sb.WriteString(fmt.Sprintf("\tm.%s.SetupRoutes(g)\n", capitalize(field.Name)))
					}
				}
			}
			sb.WriteString("}\n\n")
		}
	}

	if hasModules {
		g.project.AddFile("modules.go", sb.String())
	}
}