load("@io_bazel_rules_docker//container:container.bzl", "container_push")

container_push(
   name = "push_one",
   image = "@nginx//image",
   format = "Docker",
   registry = "{STABLE_DOCKER_REGISTRY}",
   repository = "{STABLE_REPOSITORY_NAME}",
   tag = "{STABLE_IMAGE_TAG}",
)

container_push(
   name = "push_two",
   image = "@nginx//image",
   format = "Docker",
   registry = "{STABLE_DOCKER_REGISTRY}",
   repository = "{STABLE_REPOSITORY_NAME}",
   tag = "{STABLE_IMAGE_TAG}",
)
