install-deps:
	go mod download

test: .FORCE
	go test -v ./test/

.PHONY: .FORCE
