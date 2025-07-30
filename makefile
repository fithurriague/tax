all:run 

run:
	go run cmd/main.go

build:
	go build -o dist/server cmd/main.go

test:
	go test -v ./...

bench:
	go test -bench=. -benchmem ./...
