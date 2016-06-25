default:
	go build ./cmd/load-mold
	./load-mold testdata/person.go
