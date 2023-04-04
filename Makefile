BINARY_NAME=monsterbash-server

build:
	GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_NAME}-linux main.go
	GOARCH=amd64 GOOS=darwin go build -o ./bin/${BINARY_NAME}-darwin main.go
	GOARCH=amd64 GOOS=windows go build -o ./bin/${BINARY_NAME}-windows.exe main.go
	go build -o ./bin/${BINARY_NAME} main.go

run: build
	./bin/${BINARY_NAME}

clean:
	go clean
	rm -f ./bin/${BINARY_NAME}-linux
	rm -f ./bin/${BINARY_NAME}-darwin
	rm -f ./bin/${BINARY_NAME}-windows.exe
	rm -f ./bin/${BINARY_NAME}

test:
	go test -race -v ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all