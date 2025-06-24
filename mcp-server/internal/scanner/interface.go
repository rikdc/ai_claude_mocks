package scanner

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/kohofinancial/experiments/ai_claude_prime/mcp-server/internal/types"
)

// GoInterfaceScanner scans Go source code for interface definitions
type GoInterfaceScanner struct {
	fileSet *token.FileSet
}

// NewGoInterfaceScanner creates a new interface scanner
func NewGoInterfaceScanner() *GoInterfaceScanner {
	return &GoInterfaceScanner{
		fileSet: token.NewFileSet(),
	}
}

// ScanProject scans a Go project for interface definitions
func (s *GoInterfaceScanner) ScanProject(projectPath string) ([]types.InterfaceDefinition, error) {
	var interfaces []types.InterfaceDefinition

	// Parse all Go files in the project
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-Go files
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Skip test files for now (could be configurable)
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Skip vendor directory
		if strings.Contains(path, "vendor/") {
			return nil
		}

		// Parse the Go file
		fileInterfaces, err := s.scanFile(path)
		if err != nil {
			// Log error but continue scanning other files
			return nil
		}

		interfaces = append(interfaces, fileInterfaces...)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan project: %w", err)
	}

	return interfaces, nil
}

// scanFile scans a single Go file for interface definitions
func (s *GoInterfaceScanner) scanFile(filePath string) ([]types.InterfaceDefinition, error) {
	// Parse the Go file
	src, err := parser.ParseFile(s.fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	var interfaces []types.InterfaceDefinition

	// Walk the AST to find interface declarations
	ast.Inspect(src, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.GenDecl:
			// Check if this is a type declaration
			if node.Tok == token.TYPE {
				for _, spec := range node.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						// Check if the type is an interface
						if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
							interfaceDef := s.extractInterfaceDefinition(
								typeSpec.Name.Name,
								interfaceType,
								src.Name.Name,
								filePath,
								s.fileSet.Position(typeSpec.Pos()).Line,
								node.Doc,
							)
							interfaces = append(interfaces, interfaceDef)
						}
					}
				}
			}
		}
		return true
	})

	return interfaces, nil
}

// extractInterfaceDefinition extracts interface metadata from AST nodes
func (s *GoInterfaceScanner) extractInterfaceDefinition(
	name string,
	interfaceType *ast.InterfaceType,
	packageName string,
	filePath string,
	lineNumber int,
	docGroup *ast.CommentGroup,
) types.InterfaceDefinition {
	var methods []types.MethodSignature
	var comments []string

	// Extract documentation comments
	if docGroup != nil {
		for _, comment := range docGroup.List {
			comments = append(comments, strings.TrimPrefix(comment.Text, "//"))
		}
	}

	// Extract method signatures
	for _, method := range interfaceType.Methods.List {
		if len(method.Names) > 0 {
			methodName := method.Names[0].Name
			methodSig := s.extractMethodSignature(methodName, method.Type, method.Doc)
			methods = append(methods, methodSig)
		}
	}

	return types.InterfaceDefinition{
		Name:       name,
		Package:    packageName,
		Methods:    methods,
		FilePath:   filePath,
		LineNumber: lineNumber,
		Comments:   comments,
	}
}

// extractMethodSignature extracts method signature details
func (s *GoInterfaceScanner) extractMethodSignature(
	name string,
	methodType ast.Expr,
	docGroup *ast.CommentGroup,
) types.MethodSignature {
	var parameters []types.Parameter
	var returns []types.Parameter
	var comments []string

	// Extract method documentation
	if docGroup != nil {
		for _, comment := range docGroup.List {
			comments = append(comments, strings.TrimPrefix(comment.Text, "//"))
		}
	}

	// Extract function signature details
	if funcType, ok := methodType.(*ast.FuncType); ok {
		// Extract parameters
		if funcType.Params != nil {
			for _, param := range funcType.Params.List {
				paramType := s.typeToString(param.Type)
				if len(param.Names) > 0 {
					for _, name := range param.Names {
						parameters = append(parameters, types.Parameter{
							Name: name.Name,
							Type: paramType,
						})
					}
				} else {
					// Anonymous parameter
					parameters = append(parameters, types.Parameter{
						Name: "",
						Type: paramType,
					})
				}
			}
		}

		// Extract return values
		if funcType.Results != nil {
			for _, result := range funcType.Results.List {
				resultType := s.typeToString(result.Type)
				if len(result.Names) > 0 {
					for _, name := range result.Names {
						returns = append(returns, types.Parameter{
							Name: name.Name,
							Type: resultType,
						})
					}
				} else {
					// Anonymous return value
					returns = append(returns, types.Parameter{
						Name: "",
						Type: resultType,
					})
				}
			}
		}
	}

	return types.MethodSignature{
		Name:       name,
		Parameters: parameters,
		Returns:    returns,
		Comments:   comments,
	}
}

// typeToString converts an AST type expression to string representation
func (s *GoInterfaceScanner) typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return s.typeToString(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		return "*" + s.typeToString(t.X)
	case *ast.ArrayType:
		if t.Len == nil {
			return "[]" + s.typeToString(t.Elt)
		}
		return "[" + s.typeToString(t.Len) + "]" + s.typeToString(t.Elt)
	case *ast.MapType:
		return "map[" + s.typeToString(t.Key) + "]" + s.typeToString(t.Value)
	case *ast.ChanType:
		switch t.Dir {
		case ast.SEND:
			return "chan<- " + s.typeToString(t.Value)
		case ast.RECV:
			return "<-chan " + s.typeToString(t.Value)
		default:
			return "chan " + s.typeToString(t.Value)
		}
	case *ast.FuncType:
		return "func" // Simplified for now
	case *ast.InterfaceType:
		return "interface{}"
	default:
		return "unknown"
	}
}

// ExtractInterfaceMetadata extracts detailed metadata for a specific interface
func (s *GoInterfaceScanner) ExtractInterfaceMetadata(filePath, interfaceName string) (*types.InterfaceDefinition, error) {
	interfaces, err := s.scanFile(filePath)
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		if iface.Name == interfaceName {
			return &iface, nil
		}
	}

	return nil, fmt.Errorf("interface %s not found in file %s", interfaceName, filePath)
}

// DetectDependencies analyzes import dependencies for an interface
func (s *GoInterfaceScanner) DetectDependencies(filePath string) ([]string, error) {
	src, err := parser.ParseFile(s.fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	var dependencies []string
	for _, imp := range src.Imports {
		// Remove quotes from import path
		importPath := strings.Trim(imp.Path.Value, "\"")
		dependencies = append(dependencies, importPath)
	}

	return dependencies, nil
}