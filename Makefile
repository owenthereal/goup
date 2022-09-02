SHELL=/bin/bash -o pipefail


.PHONY: vet
vet:
	shellcheck -s dash -- install.sh
	docker \
		run \
		--rm \
		-v $(CURDIR):/app \
		-w /app \
		golangci/golangci-lint:latest \
		golangci-lint run --timeout 5m -v

.PHONY: build
build:
	go build -o bin/goup ./cmd/goup

.PHONY: ftest
ftest:
	go test -v -test.count=1 -race ./ftest/...

.PHONY: docker_ftest
docker_ftest:
	docker build --rm -t jingweno/goup-ftest -f Dockerfile.ftest . && docker rmi jingweno/goup-ftest
