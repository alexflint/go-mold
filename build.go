package mold

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

func makeSkeleton(expr ast.Expr, name, pkg string) Type {
	st := staticType{name: name, pkg: pkg}
	switch expr := expr.(type) {
	case *ast.ParenExpr:
		return makeSkeleton(expr.X, name, pkg)
	case *ast.Ident:
		return &staticAlias{st: st, expr: expr}
	case *ast.SelectorExpr:
		panic("imported symbols not implemented")
	case *ast.StarExpr:
		return &staticPtr{staticType: st, expr: expr}
	case *ast.MapType:
		return &staticMap{staticType: st, expr: expr}
	case *ast.ArrayType:
		if expr.Len == nil {
			return &staticSlice{staticType: st, expr: expr}
		} else {
			return &staticArray{staticType: st, expr: expr}
		}
	case *ast.StructType:
		return &staticStruct{staticType: st, expr: expr}
	case *ast.InterfaceType:
		return &staticInterface{staticType: st, expr: expr}
	default:
		panic(fmt.Sprintf("unexpected %T in type expression", expr))
	}
}

type builder struct {
	pkg     string
	symbols map[string]Type // user-defined types + builtins
	named   map[string]Type // user-defined types only
	unnamed []Type
}

func newBuilder(pkg string) *builder {
	b := builder{
		pkg:     pkg,
		symbols: make(map[string]Type),
		named:   make(map[string]Type),
	}

	b.symbols["bool"] = TypeOf(true)
	b.symbols["byte"] = TypeOf(byte(0))
	b.symbols["rune"] = TypeOf(rune(0))
	b.symbols["string"] = TypeOf("")

	b.symbols["int"] = TypeOf(int(0))
	b.symbols["int8"] = TypeOf(int8(0))
	b.symbols["int16"] = TypeOf(int16(0))
	b.symbols["int32"] = TypeOf(int32(0))
	b.symbols["int64"] = TypeOf(int64(0))

	b.symbols["uint"] = TypeOf(uint(0))
	b.symbols["uint8"] = TypeOf(uint8(0))
	b.symbols["uint16"] = TypeOf(uint16(0))
	b.symbols["uint32"] = TypeOf(uint32(0))
	b.symbols["uint64"] = TypeOf(uint64(0))

	b.symbols["float32"] = TypeOf(float32(0))
	b.symbols["float64"] = TypeOf(float64(0))

	b.symbols["complex64"] = TypeOf(complex64(0))
	b.symbols["complex128"] = TypeOf(complex128(0))

	var err error
	b.symbols["error"] = TypeOf(&err).Elem()
	return &b
}

func (b *builder) resolve(expr ast.Expr) Type {
	switch expr := expr.(type) {
	case *ast.ParenExpr:
		return b.resolve(expr.X)
	case *ast.Ident:
		if t, found := b.symbols[expr.Name]; found {
			return t
		}
		panic(fmt.Sprintf("unknown type: %s", expr.Name))
	default:
		t := makeSkeleton(expr, "", b.pkg)
		b.populate(t)
		b.unnamed = append(b.unnamed, t)
		return t
	}
}

func (b *builder) add(decl ast.Decl) {
	if decl, ok := decl.(*ast.GenDecl); ok && decl.Tok == token.TYPE {
		for _, spec := range decl.Specs {
			spec := spec.(*ast.TypeSpec)
			name := spec.Name.Name
			t := makeSkeleton(spec.Type, name, b.pkg)
			b.symbols[name] = t
			b.named[name] = t
		}
	}
}

func (b *builder) build() {
	for _, t := range b.named {
		b.populate(t)
	}
}

func (b *builder) populate(t Type) {
	switch t := t.(type) {
	case *staticAlias:
		b.populateAlias(t)
	case *staticPtr:
		b.populatePtr(t)
	case *staticArray:
		b.populateArray(t)
	case *staticSlice:
		b.populateSlice(t)
	case *staticMap:
		b.populateMap(t)
	case *staticStruct:
		b.populateStruct(t)
	case *staticInterface:
		b.populateInterface(t)
	default:
		panic(fmt.Sprintf("unable to populate %T", t))
	}
}

func (b *builder) populateAlias(t *staticAlias) {
	t.Type = b.resolve(t.expr)
	t.expr = nil
}

func (b *builder) populatePtr(t *staticPtr) {
	t.elem = b.resolve(t.expr.X)
	t.expr = nil
}

func (b *builder) populateArray(t *staticArray) {
	t.elem = b.resolve(t.expr.Elt)

	length, ok := t.expr.Len.(*ast.BasicLit)
	if !ok {
		panic(fmt.Sprintf("can only parse arrays with literal ints for length, got %T",
			t.expr.Len))
	}
	if length.Kind != token.INT {
		panic(fmt.Sprintf("can only parse arrays with literal ints for length, got %v",
			length.Kind))
	}

	var err error
	t.length, err = strconv.Atoi(length.Value)
	if err != nil {
		panic(fmt.Sprintf("error parsing array length: %v", err))
	}
	t.expr = nil
}

func (b *builder) populateSlice(t *staticSlice) {
	t.elem = b.resolve(t.expr.Elt)
	t.expr = nil
}

func (b *builder) populateMap(t *staticMap) {
	t.key = b.resolve(t.expr.Key)
	t.elem = b.resolve(t.expr.Value)
	t.expr = nil
}

func (b *builder) populateStruct(t *staticStruct) {
	for _, f := range t.expr.Fields.List {
		r := b.resolve(f.Type)
		var tag StructTag
		if f.Tag != nil {
			tag = StructTag(f.Tag.Value)
		}
		if f.Names == nil {
			// anonymous field
			t.fields = append(t.fields, StructField{
				Name:      r.Name(),
				Type:      r,
				Tag:       tag,
				Anonymous: true,
			})
		} else {
			// symbols field
			for _, ident := range f.Names {
				t.fields = append(t.fields, StructField{
					Name: ident.Name,
					Type: r,
					Tag:  tag,
				})
			}
		}
	}
	// TODO: populate PkgPath, Offset, Index in each StructField
	t.expr = nil
}

func (b *builder) populateInterface(t *staticInterface) {
	// TODO: extract methods
	t.expr = nil
}
