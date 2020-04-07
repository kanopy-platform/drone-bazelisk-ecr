.PHONY: docker-build
docker-build:
	docker build -t drone-bazelisk-ecr .

.PHONY: test
test:
	go test ./...
