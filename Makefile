SHELL=/bin/bash -o pipefail


.PHONY: vet
vet:
	shellcheck -s dash -- install.sh
	# somehow the volume is mounted as user 1000 :shurg:
	docker run --privileged --rm -v $$(pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run -v

.PHONY: build
build:
	go build -o bin/goup ./cmd/goup

