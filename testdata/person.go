package test

type Person struct {
	Name     string
	Age      int
	Children []*Person
	maker    personMaker
	data     interface{}
}

type PersonPtr *Person

type Index map[string]*Person

type Array [10]Person

type List []Person

type personMaker interface {
	Make() Person
}
