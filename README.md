[![GoDoc](https://godoc.org/github.com/alexflint/go-mold?status.svg)](https://godoc.org/github.com/alexflint/go-mold)
[![Build Status](https://travis-ci.org/alexflint/go-mold.svg?branch=master)](https://travis-ci.org/alexflint/go-mold)

## Load Go types from source files

```shell
go get github.com/alexflint/go-mold
```

This package loads Go types from source files and returns an interface similar to `reflect.Type`.

```go
import "github.com/alexflint/go-mold"

func main() {
	types, err := mold.LoadFile("src.go")
	if err != nil {
		log.Fatal(err)
	}

	for name, typ := range types {
		fmt.Println(name)
		fmt.Println("kind:", typ.Kind())
		if typ.Kind() == reflect.Struct {
			for i := 0; i < typ.NumField(); i++ {
				fmt.Printf("field %d: %s\n", i, typ.Field(i).Name)
			}
		}
	}
}
```
