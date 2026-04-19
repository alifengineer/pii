package parser

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

func Parse(dir string) (string, []StructInfo, error) {
	pkg, err := loadPackage(dir)
	if err != nil {
		return "", nil, err
	}
	return pkg.Name, extractStructs(pkg), nil
}

func loadPackage(dir string) (*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo,
		Dir:  dir,
	}
	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return nil, fmt.Errorf("loading package: %w", err)
	}
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no packages found in %s", dir)
	}
	if len(pkgs[0].Errors) > 0 {
		return nil, fmt.Errorf("package errors: %v", pkgs[0].Errors)
	}
	return pkgs[0], nil
}

func extractStructs(pkg *packages.Package) []StructInfo {
	var structs []StructInfo

	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				_, ok = typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				obj := pkg.TypesInfo.Defs[typeSpec.Name]
				if obj == nil {
					continue
				}
				named, ok := obj.Type().(*types.Named)
				if !ok {
					continue
				}
				underlying, ok := named.Underlying().(*types.Struct)
				if !ok {
					continue
				}

				si := StructInfo{Name: typeSpec.Name.Name}
				for field := range underlying.Fields() {
					si.Fields = append(si.Fields, classifyField(field))
				}
				structs = append(structs, si)
			}
		}
	}
	return structs
}

func classifyField(field *types.Var) FieldInfo {
	fi := FieldInfo{
		Name:     field.Name(),
		TypeName: field.Type().String(),
	}

	typ := field.Type()

	if ptr, ok := typ.(*types.Pointer); ok {
		fi.IsPtr = true
		typ = ptr.Elem()
	}

	if sl, ok := typ.(*types.Slice); ok {
		fi.IsSlice = true
		typ = sl.Elem()
	}

	if named, ok := typ.(*types.Named); ok {
		fi.ElemType = named.Obj().Name()
		if _, ok := named.Underlying().(*types.Struct); ok {
			fi.IsStruct = true
		}
		fi.HasSanitize = hasSanitizeMethod(named)
	}

	return fi
}

func hasSanitizeMethod(named *types.Named) bool {
	for m := range named.Methods() {
		if m.Name() != "Sanitize" {
			continue
		}
		sig, ok := m.Type().(*types.Signature)
		if !ok {
			continue
		}
		if sig.Params().Len() != 0 {
			continue
		}
		if sig.Results().Len() != 1 {
			continue
		}
		if types.Identical(sig.Results().At(0).Type(), named) {
			return true
		}
	}

	mset := types.NewMethodSet(types.NewPointer(named))
	for sel := range mset.Methods() {
		if sel.Obj().Name() != "Sanitize" {
			continue
		}
		sig, ok := sel.Type().(*types.Signature)
		if !ok {
			continue
		}
		if sig.Params().Len() != 0 {
			continue
		}
		if sig.Results().Len() != 1 {
			continue
		}
		if types.Identical(sig.Results().At(0).Type(), named) {
			return true
		}
	}

	return false
}
