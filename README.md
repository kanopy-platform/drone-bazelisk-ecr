# drone-bazelisk-ecr

Build Image and push to a docker registry
```
make docker-build
```

Setup a secrets file for `drone exec` in _example/secrets.env_
```
ECR_REGISTRY=
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
```

Run example locally
```
cd example/
drone exec --secret-file secrets.env
```

## Useful Links

- [Amazon ECR Docker Credential Helper](https://github.com/awslabs/amazon-ecr-credential-helper)
- [bazelisk](https://github.com/bazelbuild/bazelisk)
- [Bazel Docker Rules](https://github.com/bazelbuild/rules_docker)
