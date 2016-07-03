package mold

import (
	"fmt"
	"go/ast"
	"reflect"
)

// -- staticAlias

type staticAlias struct {
	Type
	st   staticType
	expr *ast.Ident
}

func (t *staticAlias) Name() string    { return t.st.name }
func (t *staticAlias) PkgPath() string { return t.st.PkgPath() }
func (t *staticAlias) String() string  { return t.PkgPath() }

// -- staticPtr

type staticPtr struct {
	staticType
	elem Type
	expr *ast.StarExpr
}

func (t *staticPtr) Kind() reflect.Kind { return reflect.Ptr }
func (t *staticPtr) Elem() Type         { return t.elem }
func (t *staticPtr) String() string     { return "*" + t.elem.String() }

// -- staticArray

type staticArray struct {
	staticType
	elem   Type
	length int
	expr   *ast.ArrayType
}

func (t *staticArray) Kind() reflect.Kind { return reflect.Array }
func (t *staticArray) Elem() Type         { return t.elem }
func (t *staticArray) Len() int           { return t.length }
func (t *staticArray) String() string {
	return fmt.Sprintf("[%d]%s", t.length, t.elem.String())
}

// -- staticSlice

type staticSlice struct {
	staticType
	elem Type
	expr *ast.ArrayType
}

func (t *staticSlice) Kind() reflect.Kind { return reflect.Slice }
func (t *staticSlice) Elem() Type         { return t.elem }
func (t *staticSlice) String() string     { return "[]" + t.elem.String() }

// -- staticMap

type staticMap struct {
	staticType
	key  Type
	elem Type
	expr *ast.MapType
}

func (t *staticMap) Kind() reflect.Kind { return reflect.Map }
func (t *staticMap) Key() Type          { return t.key }
func (t *staticMap) Elem() Type         { return t.elem }
func (t *staticMap) String() string {
	return fmt.Sprintf("map[%s]%s", t.key.String(), t.elem.String())
}

// -- staticStruct

type staticStruct struct {
	staticType
	fields []StructField
	expr   *ast.StructType
}

func (t *staticStruct) Kind() reflect.Kind      { return reflect.Struct }
func (t *staticStruct) NumField() int           { return len(t.fields) }
func (t *staticStruct) Field(i int) StructField { return t.fields[i] }
func (t *staticStruct) String() string          { return t.PkgPath() }

// -- staticInterface

type staticInterface struct {
	staticType
	// TODO: add fields for methods
	expr *ast.InterfaceType
}

func (t *staticInterface) Kind() reflect.Kind { return reflect.Interface }
func (t *staticInterface) String() string     { return t.PkgPath() }

// -- staticType

type staticType struct {
	name string
	pkg  string
}

func newStaticType(name, pkg string) staticType {
	return staticType{
		name: name,
		pkg:  pkg,
	}
}

func (t *staticType) common() {}

// Align returns the alignment in bytes of a value of
// this type when allocated in memory.
func (t *staticType) Align() int {
	panic("Align not implemented")
}

// FieldAlign returns the alignment in bytes of a value of
// this type when used as a field in a struct.
func (t *staticType) FieldAlign() int {
	panic("FieldAlign not implemented")
}

// Method returns the i'th method in the type's method set.
// It panics if i is not in the range [0, NumMethod()).
//
// For a non-interface type T or *T, the returned Method's Type and Func
// fields describe a function whose first argument is the receiver.
//
// For an interface type, the returned Method's Type field gives the
// method signature, without a receiver, and the Func field is nil.
func (t *staticType) Method(int) reflect.Method {
	panic("Method not implemented")
}

// MethodByName returns the method with that name in the type's
// method set and a boolean indicating if the method was found.
//
// For a non-interface type T or *T, the returned Method's Type and Func
// fields describe a function whose first argument is the receiver.
//
// For an interface type, the returned Method's Type field gives the
// method signature, without a receiver, and the Func field is nil.
func (t *staticType) MethodByName(string) (reflect.Method, bool) {
	panic("MethodbyName not implemented")
}

// NumMethod returns the number of methods in the type's method set.
func (t *staticType) NumMethod() int {
	panic("NumMethod not implemented")
}

// Name returns the type's name within its package.
// It returns an empty string for unnamed types.
func (t *staticType) Name() string {
	return t.name
}

// PkgPath returns a named type's package path, that is, the import path
// that uniquely identifies the package, such as "encoding/base64".
// If the type was predeclared (string, error) or unnamed (*T, struct{}, []int),
// the package path will be the empty string.
func (t *staticType) PkgPath() string {
	if t.name == "" {
		return ""
	}
	return t.pkg + "/" + t.name
}

// Size returns the number of bytes needed to store
// a value of the given type; it is analogous to unsafe.Sizeof.
func (t *staticType) Size() uintptr {
	panic("Size not implemented")
}

// String returns a string representation of the type.
// The string representation may use shortened package names
// (e.g., base64 instead of "encoding/base64") and is not
// guaranteed to be unique among types.  To test for equality,
// compare the Types directly.
func (t *staticType) String() string {
	panic("String not implemented")
}

// Kind returns the specific kind of this type.
func (t *staticType) Kind() reflect.Kind {
	panic("Kind not implemented")
}

// Implements reports whether the type implements the interface type u.
func (t *staticType) Implements(u Type) bool {
	panic("Implements not implemented")
}

// AssignableTo reports whether a value of the type is assignable to type u.
func (t *staticType) AssignableTo(u Type) bool {
	panic("AssignableTo not implemented")
}

// ConvertibleTo reports whether a value of the type is convertible to type u.
func (t *staticType) ConvertibleTo(u Type) bool {
	panic("ConvertibleTo not implemented")
}

// Comparable reports whether values of this type are comparable.
func (t *staticType) Comparable() bool {
	panic("Comparable not implemented")
}

// Bits returns the size of the type in bits.
// It panics if the type's Kind is not one of the
// sized or unsized Int, Uint, Float, or Complex kinds.
func (t *staticType) Bits() int {
	panic("Bits not implemented")
}

// ChanDir returns a channel type's direction.
// It panics if the type's Kind is not Chan.
func (t *staticType) ChanDir() reflect.ChanDir {
	panic("ChanDir not implemented")
}

// IsVariadic reports whether a function type's final input parameter
// is a "..." parameter.  If so, t.In(t.NumIn() - 1) returns the parameter's
// implicit actual type []T.
//
// For concreteness, if t represents func(x int, y ... float64), then
//
//	t.NumIn() == 2
//	t.In(0) is the Type for "int"
//	t.In(1) is the Type for "[]float64"
//	t.IsVariadic() == true
//
// IsVariadic panics if the type's Kind is not Func.
func (t *staticType) IsVariadic() bool {
	panic("IsVariadic not implemented")
}

// Elem returns a type's element type.
// It panics if the type's Kind is not Array, Chan, Map, Ptr, or Slice.
func (t *staticType) Elem() Type {
	panic("Elem not implemented")
}

// Field returns a struct type's i'th field.
// It panics if the type's Kind is not Struct.
// It panics if i is not in the range [0, NumField()).
func (t *staticType) Field(i int) StructField {
	panic("Field not implemented")
}

// FieldByIndex returns the nested field corresponding
// to the index sequence.  It is equivalent to calling Field
// successively for each index i.
// It panics if the type's Kind is not Struct.
func (t *staticType) FieldByIndex(index []int) StructField {
	panic("FieldByIndex not implemented")
}

// FieldByName returns the struct field with the given name
// and a boolean indicating if the field was found.
func (t *staticType) FieldByName(name string) (StructField, bool) {
	panic("FieldByName not implemented")
}

// FieldByNameFunc returns the first struct field with a name
// that satisfies the match function and a boolean indicating if
// the field was found.
func (t *staticType) FieldByNameFunc(match func(string) bool) (StructField, bool) {
	panic("FieldByNameFunc not implemented")
}

// In returns the type of a function type's i'th input parameter.
// It panics if the type's Kind is not Func.
// It panics if i is not in the range [0, NumIn()).
func (t *staticType) In(i int) Type {
	panic("In not implemented")
}

// Key returns a map type's key type.
// It panics if the type's Kind is not Map.
func (t *staticType) Key() Type {
	panic("Key not implemented")
}

// Len returns an array type's length.
// It panics if the type's Kind is not Array.
func (t *staticType) Len() int {
	panic("Len not implemented")
}

// NumField returns a struct type's field count.
// It panics if the type's Kind is not Struct.
func (t *staticType) NumField() int {
	panic("NumField not implemented")
}

// NumIn returns a function type's input parameter count.
// It panics if the type's Kind is not Func.
func (t *staticType) NumIn() int {
	panic("NumIn not implemented")
}

// NumOut returns a function type's output parameter count.
// It panics if the type's Kind is not Func.
func (t *staticType) NumOut() int {
	panic("NumOut not implemented")
}

// Out returns the type of a function type's i'th output parameter.
// It panics if the type's Kind is not Func.
// It panics if i is not in the range [0, NumOut()).
func (t *staticType) Out(i int) Type {
	panic("Out not implemented")
}
