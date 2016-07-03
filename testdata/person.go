package test

type Person struct {
	Name     string
	Age      int
	Children []PersonPtr
	maker    personMaker
	data     interface{}
}

type PersonPtr *Person

type PersonArray [10]Person

type PersonSlice []Person

type AddressBook map[string]PersonPtr

type personMaker interface {
	Make() Person
}
