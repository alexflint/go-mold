package mold

import "reflect"

type liveType struct {
	reflect.Type
}

func (t liveType) AssignableTo(u Type) bool {
	if u, ok := u.(liveType); ok {
		return t.Type.AssignableTo(u.Type)
	}
	return false
}

func (t liveType) ConvertibleTo(u Type) bool {
	if u, ok := u.(liveType); ok {
		return t.Type.ConvertibleTo(u.Type)
	}
	return false
}

func (t liveType) Implements(u Type) bool {
	if u, ok := u.(liveType); ok {
		return t.Type.Implements(u.Type)
	}
	return false
}

func (t liveType) Key() Type {
	return liveType{t.Type.Key()}
}

func (t liveType) Elem() Type {
	return liveType{t.Type.Elem()}
}

func (t liveType) In(i int) Type {
	return liveType{t.Type.In(i)}
}

func (t liveType) Out(i int) Type {
	return liveType{t.Type.Out(i)}
}

func structField(f reflect.StructField) StructField {
	return StructField{
		Name:      f.Name,
		PkgPath:   f.PkgPath,
		Type:      liveType{f.Type},
		Tag:       StructTag(f.Tag),
		Offset:    f.Offset,
		Anonymous: f.Anonymous,
	}
}

func (t liveType) Field(i int) StructField {
	return structField(t.Type.Field(i))
}

func (t liveType) FieldByIndex(i []int) StructField {
	return structField(t.Type.FieldByIndex(i))
}

func (t liveType) FieldByName(name string) (StructField, bool) {
	f, b := t.Type.FieldByName(name)
	return structField(f), b
}

func (t liveType) FieldByNameFunc(g func(string) bool) (StructField, bool) {
	f, b := t.Type.FieldByNameFunc(g)
	return structField(f), b
}

// construct a type from a live object
func typeOf(v interface{}) Type {
	return liveType{reflect.TypeOf(v)}
}
