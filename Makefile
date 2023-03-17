.PHONY: docker-build
docker-build:
	docker build -t registry.example.com/drone-bazelisk-ecr .

.PHONY: test
test:
	go vet ./...
	go test ./...
