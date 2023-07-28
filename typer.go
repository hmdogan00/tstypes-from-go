package main

import (
	"fmt"
	"go/ast"
)

func getTypeNameFromNode(node ast.Node) string {
	if typeIdent, ok := node.(*ast.Ident); ok {
		return typeIdent.Name
	}
	if typeSelector, ok := node.(*ast.SelectorExpr); ok {
		return typeSelector.Sel.Name
	}
	if typeStarExpr, ok := node.(*ast.StarExpr); ok {
		return getTypeNameFromNode(typeStarExpr.X)
	}
	if typeArrayType, ok := node.(*ast.ArrayType); ok {
		return getTypeNameFromNode(typeArrayType.Elt)
	}
	if typeMapType, ok := node.(*ast.MapType); ok {
		return getTypeNameFromNode(typeMapType.Value)
	}
	if basicVal, ok := node.(*ast.BasicLit); ok {
		return basicVal.Kind.String()
	}
	if compositeVal, ok := node.(*ast.CompositeLit); ok {
		if arrayType, ok := compositeVal.Type.(*ast.ArrayType); ok {
			if arrayTypeIdent, ok := arrayType.Elt.(*ast.Ident); ok {
				return arrayTypeIdent.Name
			} else {
				return getTypeNameFromNode(arrayType.Elt)
			}
		} else {
			return getTypeNameFromNode(compositeVal.Type)
		}
	}
	panic(fmt.Sprintf("Unknown type: %T", node))
}

func varHandler(valueSpec ast.ValueSpec) string {
	result := ""
	for _, name := range valueSpec.Names {
		typeIdent, ok := valueSpec.Type.(*ast.Ident)
		// basic cases of var foo int32 = 32
		if ok {
			result += fmt.Sprintf("export type %s = %s;\n", name.Name, getTsTypeName(typeIdent.Name))
		} else {
			// case of var foo = 32
			if basicVal, ok := valueSpec.Values[0].(*ast.BasicLit); ok {
				result += fmt.Sprintf("export type %s = %s;\n", name.Name, getTsTypeName(basicVal.Kind.String()))
			}
			if compositeVal, ok := valueSpec.Values[0].(*ast.CompositeLit); ok {
				if arrayType, ok := compositeVal.Type.(*ast.ArrayType); ok {
					if arrayTypeIdent, ok := arrayType.Elt.(*ast.Ident); ok {
						result += fmt.Sprintf("export type %s = %s[];\n", name.Name, getTsTypeName(arrayTypeIdent.Name))
					}
				}
			}
		}
	}
	return result
}

func interfaceHandler(typeSpec ast.TypeSpec, structType ast.StructType) string {
	result := fmt.Sprintf("export interface %s {\n", typeSpec.Name.Name)
	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			tsName := getTsTypeName(getTypeNameFromNode(field.Type))
			result += fmt.Sprintf("\t%s: %s,\n", name.Name, getTsTypeName(tsName))
		}
	}
	return result + "}\n"
}

func functionHandler(funcDecl ast.FuncDecl) string {
	result := fmt.Sprintf("export function %s(", funcDecl.Name.Name)
	listLength := len(funcDecl.Type.Params.List)
	count := 0
	for _, field := range funcDecl.Type.Params.List {
		for _, name := range field.Names {
			tsName := getTsTypeName(getTypeNameFromNode(field.Type))
			count++
			format := "%s: %s, "
			if count == listLength {
				format := "%s: %s"
				result += fmt.Sprintf(format, name.Name, tsName)
				continue
			}
			result += fmt.Sprintf(format, name.Name, tsName)
		}
	}
	result += "): "
	if funcDecl.Type.Results != nil {
		listLength := len(funcDecl.Type.Results.List)
		count := 0
		for _, field := range funcDecl.Type.Results.List {
			tsName := getTsTypeName(getTypeNameFromNode(field.Type))
			count++
			format := "%s, "
			if count == listLength {
				format = "%s"
				result += fmt.Sprintf(format, tsName)
				continue
			}
			result += fmt.Sprintf("%s, ", tsName)
		}
	} else {
		result += "void"
	}
	return result + ";\n"
}

func enumHandler(specs []ast.Spec) string {
	result := "export enum {\n"
	for _, spec := range specs {
		if valueSpec, ok := spec.(*ast.ValueSpec); ok {
			for _, name := range valueSpec.Names {
				result += fmt.Sprintf("\t%s = ", name.Name)
			}
			for _, value := range valueSpec.Values {
				if basicVal, ok := value.(*ast.BasicLit); ok {
					result += fmt.Sprintf("%s,\n", basicVal.Value)
				} else {
					fmt.Println("Unhandled const type detected")
					panic(value)
				}
			}
		} else {
			fmt.Println("Unhandled const type detected")
			panic(spec)
		}
	}
	result += "}\n"
	return result
}
