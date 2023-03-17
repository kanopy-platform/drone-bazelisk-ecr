uname:=$(shell uname -m)

ifeq "$(uname)" "x86_64"
	ARCH := amd64
else
	ARCH := $(uname)
endif

.PHONY: docker-build
docker-build:
	docker build --build-arg ARCH=$(ARCH) -t registry.example.com/drone-bazelisk-ecr .

.PHONY: test
test:
	go vet ./...
	go test ./...
