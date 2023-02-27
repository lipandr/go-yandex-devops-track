build:
	go build -o main *.go

#test:
#	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

vet:
	go vet

lint:
	golangci-lint run --enable-all
