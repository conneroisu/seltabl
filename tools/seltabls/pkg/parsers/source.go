package parsers

import (
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"unicode"

	"golang.org/x/tools/go/packages"
)

// GoFile is a struct for a go file
//
// https://golang.org/pkg/go/ast/#File
type GoFile struct {
	Package         string
	Path            string
	GlobalConstants []*GoType
	GlobalVariables []*GoType
	Structs         []*GoStruct
	Interfaces      []*GoInterface
	Imports         []*GoImport
	StructMethods   []*GoStructMethod
}

// PackImporter is a struct for a package importer
//
// https://golang.org/pkg/go/types/#Importer
type PackImporter struct {
	Fset *token.FileSet
}

// Import imports a package
//
// https://golang.org/pkg/go/types/#Importer
func (this *PackImporter) Import(path string) (*types.Package, error) {
	println("searching for " + path)
	pack, err := importer.Default().Import(path)
	if err != nil {
		cfg := &packages.Config{
			Fset:  this.Fset,
			Mode:  packages.NeedTypes,
			Tests: true,
		}
		// Load the package by its import path
		pkgs, err := packages.Load(cfg, path)
		if err != nil {
			return nil, err
		}
		// Check for errors
		if packages.PrintErrors(pkgs) > 0 {
			return nil, fmt.Errorf("package %s has errors", path)
		}
		// Return the first package object
		return pkgs[0].Types, nil
	}
	return pack, nil
}

// ImportPath returns the import path for the go file
func (g *GoFile) ImportPath() (
	importPath string,
	isExternalPackage bool,
	err error,
) {
	isExternalPackage = false
	importPath, err = filepath.Abs(g.Path)
	if err != nil {
		return "", false, err
	}
	if _, err = os.Stat(importPath); err != nil {
		return g.Path, false, err
	}
	if !isInGoPackages(importPath) {
		importPath = strings.TrimSuffix(importPath, filepath.Base(importPath))
		importPath = strings.TrimSuffix(importPath, "/")
		return importPath, false, nil
	}
	importPath, err = filepath.Abs(g.Path)
	if err != nil {
		return
	}
	importPath = strings.Replace(importPath, "\\", "/", -1)
	goPath := strings.Replace(build.Default.GOPATH, "\\", "/", -1)
	isExternalPackage = true
	importPath = strings.TrimPrefix(importPath, goPath)
	importPath = strings.TrimPrefix(importPath, "/src/")
	importPath = strings.TrimPrefix(importPath, "/pkg/mod/")
	importPath = strings.TrimPrefix(importPath, "pkg/mod/")
	i := strings.Index(importPath, "@")
	if i > 0 {
		importPath = importPath[:i]
	}
	if strings.HasSuffix(strings.ToLower(importPath), ".go") {
		i := strings.LastIndex(importPath, "/")
		if i > 0 {
			importPath = importPath[:i]
		}
	}
	if strings.Contains(importPath, "!") { // replace "!c" to "C"
		temp := ""
		nextUppercase := false
		for i := 0; i < len(importPath); i++ {
			if importPath[i] == '!' {
				nextUppercase = true
			} else {
				if nextUppercase {
					temp += strings.ToUpper(string(importPath[i]))
					nextUppercase = false
				} else {
					temp += string(importPath[i])
				}
			}
		}
		importPath = temp
	}
	importPath = strings.TrimSuffix(importPath, "/")
	return
}

// GoImport is a struct for an import within a file
//
// https://golang.org/pkg/go/ast/#ImportSpec
type GoImport struct {
	File *GoFile
	Name string
	Path string
}

// GoInterface is a struct for an interface within a file
//
// https://golang.org/pkg/go/ast/#TypeSpec
type GoInterface struct {
	File     *GoFile
	Name     string
	Comments string
	Methods  []*GoMethod
}

// GoMethod is a struct for a method within a file
// https://golang.org/pkg/go/ast/#FuncDecl
type GoMethod struct {
	Name     string
	Params   []*GoType
	Comments string
	Results  []*GoType
}

// GoStructMethod is a struct for a struct method within a file.
//
// https://golang.org/pkg/go/ast/#FuncDecl
type GoStructMethod struct {
	GoMethod
	Receivers []string
}

// GoType is a struct for a type within a file.
//
// https://golang.org/pkg/go/ast/#Expr
type GoType struct {
	Name       string
	Type       string
	Underlying string
	Inner      []*GoType
}

// GoStruct is a struct for a struct within a file.
//
// https://golang.org/pkg/go/ast/#TypeSpec
type GoStruct struct {
	File     *GoFile
	StartRow int
	EndRow   int
	Name     string
	Comments string
	Fields   []*GoField

	SeltablURL        string
	SeltablOccurances int
	SeltablIgnores    []string
}

// GoField is a struct for a field within a struct
//
// https://golang.org/pkg/go/ast/#Field
type GoField struct {
	Struct *GoStruct
	Name   string
	Type   string
	Tag    *GoTag
}
type Dimensions struct {
	RowStart    int
	RowEnd      int
	ColumnStart int
	ColumnEnd   int
}

// GoTag is a struct for a tag within a field
//
// https://golang.org/pkg/go/ast/#Field
type GoTag struct {
	Field *GoField
	Value string
}

// Get returns the value of the tag for the given key
func (g *GoTag) Get(key string) string {
	tag := strings.Replace(g.Value, "`", "", -1)
	return reflect.StructTag(tag).Get(key)
}

// For an import - guess what prefix will be used
// in type declarations.  For examples:
//
//	"strings" -> "strings"
//	"net/http/httptest" -> "httptest"
//
// Libraries where the package name does not match
// will be mis-identified.
func (g *GoImport) Prefix() string {
	if g.Name != "" {
		return g.Name
	}
	path := strings.Trim(g.Path, "\"")
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == -1 {
		return path
	}
	return path[lastSlash+1:]
}

// ParseFiles parses files at the same time
func ParseDir(
	path string,
	withComments bool,
	filterFiles func(fs.FileInfo) bool,
) ([]*GoFile, error) {
	// File: A File node represents a Go source file: https://golang.org/pkg/go/ast/#File
	fset := token.NewFileSet()
	var mode parser.Mode
	if withComments {
		mode = parser.ParseComments
	} else {
		mode = 0
	}
	pkgs, err := parser.ParseDir(fset, path, filterFiles, mode)
	if err != nil {
		return nil, err
	}
	goFiles := make([]*GoFile, 0)
	for _, astPackage := range pkgs {
		files := make([]*ast.File, 0)
		for _, file := range astPackage.Files {
			files = append(files, file)
		}
		for name, file := range astPackage.Files {
			goFile, err := parseFile(name, file, fset, files)
			if err != nil {
				return nil, err
			}
			goFiles = append(goFiles, goFile)
		}
	}
	return goFiles, nil
}

// ParseSingleFile parses a single file at the same time
func ParseSingleFile(path string, withComments bool) (*GoFile, error) {
	fset := token.NewFileSet()
	var mode parser.Mode
	if withComments {
		mode = parser.ParseComments
	} else {
		mode = 0
	}
	file, err := parser.ParseFile(fset, path, nil, mode)
	if err != nil {
		return nil, err
	}
	return parseFile(path, file, fset, []*ast.File{file})
}

// ParseSource parses the source of a given golang file.
func ParseSource(source string, filepath string, withComments bool) (*GoFile, error) {
	fset := token.NewFileSet()
	path := filepath
	var mode parser.Mode
	if withComments {
		mode = parser.ParseComments
	} else {
		mode = 0
	}
	file, err := parser.ParseFile(fset, path, source, mode)
	if err != nil {
		return nil, err
	}
	return parseFile(path, file, fset, []*ast.File{file})
}
func execCommand(name string, args ...string) (out string, exitCode int, err error) {
	stream := &strings.Builder{}
	cmd := exec.Command(name, args...)
	cmd.Stderr = stream
	cmd.Stdout = stream
	fmt.Printf("%v\n", strings.Join(cmd.Args, " "))
	err = cmd.Run()
	if err != nil {
		var terr *exec.ExitError
		if errors.As(err, &terr) {
			exitCode = terr.ExitCode()
			out = stream.String()
		}
	}
	fmt.Printf("Execution: %v\n", stream.String())
	return
}
func getLibrary(importUrl string) (err error, cleanup func()) {
	fmt.Printf("Importing %v\n", importUrl)
	cleanup = func() {}
	var out string
	var exitCode int
	_, staterr := os.Stat("go.mod")
	if os.IsNotExist(staterr) {
		out, exitCode, err = execCommand("go", "mod", "init", "tempmod")
		if err != nil {
			err = fmt.Errorf("failed to execute go mod init command to import Go library: %v.\nError: %v. Exit Code: %v\nOutput: %v\n", importUrl, err, exitCode, out)
			return
		}
		cleanup = func() {
			_ = os.Remove("go.mod")
			_ = os.Remove("go.sum")
		}
	}
	out, exitCode, err = execCommand("go", "get", "-v", importUrl)
	if err != nil {
		err = fmt.Errorf("failed to execute go get command to import Go library: %v.\nError: %v. Exit Code: %v\nOutput: %v\n", importUrl, err, exitCode, out)
		return
	}
	return
}
func isNamePublic(name string) bool {
	if name == "" {
		return false
	}
	r := []rune(name)[0]
	return unicode.IsLetter(r) && unicode.IsUpper(r)
}

// parseFile parses a file
func parseFile(
	path string,
	file *ast.File,
	fset *token.FileSet,
	files []*ast.File,
) (*GoFile, error) {
	var err error
	// To import sources from vendor, we use "source" compile
	// https://github.com/golang/go/issues/11415#issuecomment-283445198
	conf := types.Config{Importer: &PackImporter{fset} /*importer.ForCompiler(fset, "source", nil)*/}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
	conf.IgnoreFuncBodies = true
	tries := 2
	for tries > 0 {
		tries--
		if _, err = conf.Check(file.Name.Name, fset, files, info); err != nil {
			// Get package to import
			startingPointString := "could not import "
			start := strings.Index(err.Error(), startingPointString)
			if start > -1 && tries > 0 {
				start += len(startingPointString)
				end := strings.Index(err.Error()[start:], " ")
				if end > -1 {
					toimport := err.Error()[start : start+end]
					err, cleanup := getLibrary(toimport)
					defer cleanup()
					if err != nil {
						return nil, err
					}
					continue
				}
			}
			return nil, fmt.Errorf("errors type checking source file. error: %v", err)
		}
	}
	goFile := &GoFile{
		Path:    path,
		Package: file.Name.Name,
		Structs: []*GoStruct{},
	}
	// File.Decls: A list of the declarations in the file: https://golang.org/pkg/go/ast/#Decl
	for _, decl := range file.Decls {
		switch declType := decl.(type) {
		// GenDecl: represents an import, constant, type or variable declaration: https://golang.org/pkg/go/ast/#GenDecl
		case *ast.GenDecl:
			genDecl := declType
			// Specs: the Spec type stands for any of *ImportSpec, *ValueSpec, and *TypeSpec: https://golang.org/pkg/go/ast/#Spec
			for _, genSpec := range genDecl.Specs {
				switch genSpecType := genSpec.(type) {
				// A ValueSpec node represents a constant or variable declaration: https://pkg.go.dev/go/ast#ValueSpec
				case *ast.ValueSpec:
					if !isNamePublic(genSpecType.Names[0].Name) {
						continue
					}
					valueSpec := genSpecType
					switch genDecl.Tok {
					case token.CONST:
						goConst := buildGoConstant(goFile, info, valueSpec)
						goFile.GlobalConstants = append(goFile.GlobalConstants, goConst)
					case token.VAR:
						goVar := buildGoVariable(goFile, info, valueSpec)
						goFile.GlobalVariables = append(goFile.GlobalVariables, goVar)
					}
				// TypeSpec: A TypeSpec node represents a type declaration: https://golang.org/pkg/go/ast/#TypeSpec
				case *ast.TypeSpec:
					if !isNamePublic(genSpecType.Name.Name) {
						continue
					}
					typeSpec := genSpecType
					// typeSpec.Type: an Expr (expression) node: https://golang.org/pkg/go/ast/#Expr
					switch typeSpecType := typeSpec.Type.(type) {
					// StructType: A StructType node represents a struct type: https://golang.org/pkg/go/ast/#StructType
					case (*ast.StructType):
						structType := typeSpecType
						goStruct, err := buildGoStruct(goFile, info, typeSpec, structType, genDecl.Doc)
						if err != nil {
							return nil, err
						}
						// add the start and end row to the struct
						pos := fset.Position(structType.Pos())
						// TODO: May be off by one error here
						goStruct.StartRow = pos.Line
						goStruct.EndRow = pos.Line + structType.Fields.NumFields() + 1
						goFile.Structs = append(goFile.Structs, goStruct)
					// InterfaceType: An InterfaceType node represents an interface type. https://golang.org/pkg/go/ast/#InterfaceType
					case (*ast.InterfaceType):
						interfaceType := typeSpecType
						goInterface := buildGoInterface(goFile, info, typeSpec, interfaceType, genDecl.Doc)
						goFile.Interfaces = append(goFile.Interfaces, goInterface)
					default:
						// a not-implemented typeSpec.Type.(type), ignore
					}
					// ImportSpec: An ImportSpec node represents a single package import. https://golang.org/pkg/go/ast/#ImportSpec
				case *ast.ImportSpec:
					importSpec := genSpec.(*ast.ImportSpec)
					goImport := buildGoImport(importSpec, goFile)
					goFile.Imports = append(goFile.Imports, goImport)
				default:
					// a not-implemented genSpec.(type), ignore
				}
			}
		case *ast.FuncDecl:
			if !isNamePublic(declType.Name.Name) {
				continue
			}
			funcDecl := declType
			goStructMethod := buildStructMethod(info, funcDecl, declType.Doc)
			goFile.StructMethods = append(goFile.StructMethods, goStructMethod)
		default:
			// a not-implemented decl.(type), ignore
		}
	}
	return goFile, nil
}
func buildGoVariable(_ *GoFile, info *types.Info, spec *ast.ValueSpec) *GoType {
	var t *GoType
	if spec.Type == nil { // untyped const
		t = buildType(info, spec.Values[0])
	} else {
		t = buildType(info, spec.Type)
	}
	t.Name = spec.Names[0].Name
	return t
}
func goTypeStringFromInterface(data interface{}) string {
	return reflect.TypeOf(data).Name()
}
func buildGoConstant(_ *GoFile, info *types.Info, spec *ast.ValueSpec) *GoType {
	var t *GoType
	if spec.Type != nil { // untyped const
		t = buildType(info, spec.Type)
	} else if len(spec.Values) > 0 {
		t = buildType(info, spec.Values[0])
	} else {
		t = &GoType{Type: goTypeStringFromInterface(spec.Names[0].Obj.Data)}
	}
	t.Name = spec.Names[0].Name
	return t
}
func buildGoImport(spec *ast.ImportSpec, file *GoFile) *GoImport {
	name := ""
	if spec.Name != nil {
		name = spec.Name.Name
	}
	path := ""
	if spec.Path != nil {
		path = spec.Path.Value
	}
	return &GoImport{
		Name: name,
		Path: path,
		File: file,
	}
}
func buildGoInterface(
	file *GoFile,
	info *types.Info,
	typeSpec *ast.TypeSpec,
	interfaceType *ast.InterfaceType,
	cg *ast.CommentGroup,
) *GoInterface {
	goInterface := &GoInterface{
		File:     file,
		Name:     typeSpec.Name.Name,
		Methods:  buildMethodList(info, interfaceType.Methods.List),
		Comments: extractComment(cg),
	}
	return goInterface
}

// buildMethodList builds a list of GoMethod from the given field list
func buildMethodList(
	info *types.Info,
	fieldList []*ast.Field,
) []*GoMethod {
	methods := []*GoMethod{}
	for _, field := range fieldList {
		name := getNames(field)[0]
		fType, ok := field.Type.(*ast.FuncType)
		if !ok {
			// method was not a function
			continue
		}
		goMethod := &GoMethod{
			Name:    name,
			Params:  buildTypeList(info, fType.Params),
			Results: buildTypeList(info, fType.Results),
		}
		methods = append(methods, goMethod)
	}
	return methods
}

// buildStructMethod builds a GoStructMethod from the given function declaration and comment group
func buildStructMethod(
	info *types.Info,
	funcDecl *ast.FuncDecl,
	cg *ast.CommentGroup,
) *GoStructMethod {
	return &GoStructMethod{
		Receivers: buildReceiverList(info, funcDecl.Recv),
		GoMethod: GoMethod{
			Name:     funcDecl.Name.Name,
			Params:   buildTypeList(info, funcDecl.Type.Params),
			Results:  buildTypeList(info, funcDecl.Type.Results),
			Comments: extractComment(cg),
		},
	}
}

// buildReceiverList builds a list of receivers from the given field list
func buildReceiverList(info *types.Info, fieldList *ast.FieldList) []string {
	receivers := []string{}
	if fieldList != nil {
		for _, t := range fieldList.List {
			receivers = append(receivers, getTypeString(info, t.Type))
		}
	}
	return receivers
}

// buildTypeList builds a list of GoType from the given field list
func buildTypeList(info *types.Info, fieldList *ast.FieldList) []*GoType {
	types := []*GoType{}
	if fieldList != nil {
		for _, t := range fieldList.List {
			goType := buildType(info, t.Type)
			for _, n := range getNames(t) {
				copyType := copyType(goType)
				copyType.Name = n
				types = append(types, copyType)
			}
		}
	}
	return types
}

// getNames gets the names of the given field
func getNames(field *ast.Field) []string {
	if len(field.Names) == 0 {
		return []string{""}
	}
	result := []string{}
	for _, name := range field.Names {
		result = append(result, name.String())
	}
	return result
}

// getTypeString gets the type string of the given type expression
func getTypeString(info *types.Info, expr ast.Expr) string {
	if expr == nil {
		return ""
	}
	if typeInfo := info.TypeOf(expr); typeInfo != nil {
		return typeInfo.String()
	}
	panic("info.TypeOf failed to extract type from expression")
}

// typesBasicToString converts a types.Basic to a string
func typesBasicToString(basic *types.Basic) string {
	switch basic.Kind() {
	case types.Bool:
		return "bool"
	case types.Int:
		return "int"
	case types.Int8:
		return "int8"
	case types.Int16:
		return "int16"
	case types.Int32:
		return "int32"
	case types.Int64:
		return "int64"
	case types.Uint:
		return "uint"
	case types.Uint8:
		return "uint8"
	case types.Uint16:
		return "uint16"
	case types.Uint32:
		return "uint32"
	case types.Uint64:
		return "uint64"
	case types.Uintptr:
		return "uint64"
	case types.Float32:
		return "float32"
	case types.Float64:
		return "float64"
	case types.Complex64:
		return "complex64"
	case types.Complex128:
		return "complex128"
	case types.String:
		return "string"
	case types.UnsafePointer:
		return "uint64"
		// types for untyped values
	case types.UntypedBool:
		return "bool"
	case types.UntypedInt:
		return "int"
	case types.UntypedRune:
		return "int32"
	case types.UntypedFloat:
		return "float64"
	case types.UntypedComplex:
		return "complex128"
	case types.UntypedString:
		return "string"
	case types.UntypedNil:
		return ""
	default:
		return ""
	}
}

// getUnderlyingTypeString gets the underlying type string of the given type expression
func getUnderlyingTypeString(info *types.Info, expr ast.Expr) string {
	if typeInfo := info.TypeOf(expr); typeInfo != nil {
		if underlying := typeInfo.Underlying(); underlying != nil {
			if _, ok := underlying.(*types.Interface); ok {
				return typeInfo.String()
			}
			if _, ok := underlying.(*types.Slice); ok {
				if e, ok := underlying.(*types.Slice).Elem().(*types.Basic); ok {
					str := typesBasicToString(e)
					if str != "" {
						return "[]" + str
					}
				}
			}
			if e, ok := underlying.(*types.Basic); ok {
				str := typesBasicToString(e)
				if str != "" {
					return str
				}
			}
			return underlying.String()
		}
	}
	return ""
}

// copyType copies a GoType
func copyType(goType *GoType) *GoType {
	return &GoType{
		Type:       goType.Type,
		Inner:      goType.Inner,
		Name:       goType.Name,
		Underlying: goType.Underlying,
	}
}

// buildType builds a GoType from the given type expression
func buildType(info *types.Info, expr ast.Expr) *GoType {
	innerTypes := []*GoType{}
	typeString := getTypeString(info, expr)
	underlyingString := getUnderlyingTypeString(info, expr)
	switch specType := expr.(type) {
	case *ast.FuncType:
		innerTypes = append(innerTypes, buildTypeList(info, specType.Params)...)
		innerTypes = append(innerTypes, buildTypeList(info, specType.Results)...)
	case *ast.ArrayType:
		innerTypes = append(innerTypes, buildType(info, specType.Elt))
	case *ast.StructType:
		innerTypes = append(innerTypes, buildTypeList(info, specType.Fields)...)
	case *ast.MapType:
		innerTypes = append(innerTypes, buildType(info, specType.Key))
		innerTypes = append(innerTypes, buildType(info, specType.Value))
	case *ast.ChanType:
		innerTypes = append(innerTypes, buildType(info, specType.Value))
	case *ast.StarExpr:
		innerTypes = append(innerTypes, buildType(info, specType.X))
	case *ast.Ellipsis:
		typeString = strings.ReplaceAll(typeString, "[]", "...")
		underlyingString = strings.ReplaceAll(underlyingString, "[]", "...")
		innerTypes = append(innerTypes, buildType(info, specType.Elt))
	case *ast.InterfaceType:
		methods := buildMethodList(info, specType.Methods.List)
		for _, m := range methods {
			innerTypes = append(innerTypes, m.Params...)
			innerTypes = append(innerTypes, m.Results...)
		}
	case *ast.Ident:
	case *ast.SelectorExpr:
	case *ast.BasicLit:
	case *ast.BinaryExpr:
	default:
		fmt.Printf("Unexpected field type: `%s`,\n %#v\n", typeString, specType)
	}
	return &GoType{
		Type:       typeString,
		Underlying: underlyingString,
		Inner:      innerTypes,
	}
}

// buildGoStruct builds a GoStruct from the given type spec and comment group
func buildGoStruct(
	file *GoFile,
	info *types.Info,
	typeSpec *ast.TypeSpec,
	structType *ast.StructType,
	cg *ast.CommentGroup,
) (*GoStruct, error) {
	comment := extractComment(cg)
	url, ignores, occurances, err := parseStructComments(comment)
	if err != nil {
		return nil, err
	}
	goStruct := &GoStruct{
		File:              file,
		Name:              typeSpec.Name.Name,
		Fields:            []*GoField{},
		Comments:          comment,
		SeltablURL:        url,
		SeltablOccurances: occurances,
		SeltablIgnores:    ignores,
	}
	// Field: A Field declaration list in a struct type, a method list in an interface type,
	// or a parameter/result declaration in a signature: https://golang.org/pkg/go/ast/#Field
	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			goField := &GoField{
				Struct: goStruct,
				Name:   name.String(),
				Type:   getTypeString(info, field.Type),
			}
			if field.Tag != nil {
				goTag := &GoTag{
					Field: goField,
					Value: field.Tag.Value,
				}
				goField.Tag = goTag
			}
			goStruct.Fields = append(goStruct.Fields, goField)
		}
	}
	return goStruct, nil
}

// extractComment extracts the comment from the given comment group
func extractComment(cg *ast.CommentGroup) string {
	if cg == nil || cg.List == nil {
		return ""
	}
	var comment string
	for _, c := range cg.List {
		comment += c.Text
		comment = strings.ReplaceAll(comment, "//", "")
		comment = strings.ReplaceAll(comment, "/*", "")
		comment = strings.ReplaceAll(comment, "*/", "")
		comment = strings.TrimSpace(comment)
	}
	return comment
}

// isInGoPackages checks if the given path is in the go packages
func isInGoPackages(path string) bool {
	goPath := strings.Replace(build.Default.GOPATH, "\\", "/", -1)
	return strings.Contains(path, goPath)
}
