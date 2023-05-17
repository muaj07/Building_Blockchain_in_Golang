build:
	go build -o ./bin/transport
run: build
	./bin/transport
test:
	go test -v ./...
	