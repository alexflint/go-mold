package mold

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFile_Person(t *testing.T) {
	types, err := LoadFile("testdata/person.go")
	require.NoError(t, err)

	person := types["Person"]
	personPtr := types["PersonPtr"]
	personArray := types["PersonArray"]
	personSlice := types["PersonSlice"]
	personMaker := types["personMaker"]
	addressBook := types["AddressBook"]

	require.NotNil(t, person)
	require.NotNil(t, personPtr)
	require.NotNil(t, personArray)
	require.NotNil(t, personSlice)
	require.NotNil(t, personMaker)
	require.NotNil(t, addressBook)

	assert.Equal(t, reflect.Struct, person.Kind())
	require.Equal(t, 5, person.NumField())

	nameField := person.Field(0)
	assert.Equal(t, "Name", nameField.Name)

	ageField := person.Field(1)
	assert.Equal(t, "Age", ageField.Name)

	childrenField := person.Field(2)
	assert.Equal(t, "Children", childrenField.Name)

	makerField := person.Field(3)
	assert.Equal(t, "maker", makerField.Name)
}
