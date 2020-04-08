# drone-bazelisk-ecr

Drone plugin for building images with Bazel rules_docker and ECR.

This plugin sets the following environment variables during builds so that they can be referenced as stamp variables in workspace status scripts.

    DRONE_ECR_REGISTRY
    DRONE_ECR_REPOSITORY
    DRONE_ECR_TAG

See the [example directory](./example) to see how this plugin interacts with your build environment.

## Testing locally with `drone exec`

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
