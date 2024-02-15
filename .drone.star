# workaround to render locally since you cant pass repo.branch to the cli
def repo_branch(ctx):
    return getattr(ctx.repo, "branch", "main")


def version(ctx):
    # use git commit if this is not a tag event
    if ctx.build.event != "tag":
        return "git-{}".format(commit(ctx))

    return ctx.build.ref.removeprefix("refs/tags/")


def version_tag(ctx, arch):
    return "{}-{}".format(version(ctx), arch)


def commit(ctx):
    return ctx.build.commit[:7]


def build_env(ctx):
    return {
        "GIT_COMMIT": commit(ctx),
        "VERSION": version(ctx),
    }


def new_pipeline(name, arch, **kwargs):
    pipeline = {
        "kind": "pipeline",
        "name": name,
        "platform": {
            "arch": arch,
        },
        "steps": [],
    }

    pipeline.update(kwargs)

    return pipeline


def pipeline_test(ctx):
    cache_volume = {"name": "cache", "temp": {}}
    cache_mount = {"name": "cache", "path": "/go"}

    # licensed-go image only supports amd64
    return new_pipeline(
        name="test",
        arch="amd64",
        trigger={"branch": repo_branch(ctx)},
        volumes=[cache_volume],
        workspace={"path": "/go/src/github.com/{}".format(ctx.repo.slug)},
        steps=[
            {
                "commands": ["make test"],
                "image": "golangci/golangci-lint",
                "name": "test",
                "volumes": [cache_mount],
            },
            {
                "commands": ["licensed cache", "licensed status"],
                "image": "public.ecr.aws/kanopy/licensed-go",
                "name": "license-check",
            },
            {
                "image": "plugins/kaniko-ecr",
                "name": "build",
                "pull": "always",
                "settings": {"no_push": True},
                "volumes": [cache_mount],
                "when": {"event": ["pull_request"]},
            },
        ],
    )


def pipeline_build(ctx, arch):
    return new_pipeline(
        depends_on=["test"],
        name="publish-{}".format(arch),
        arch=arch,
        steps=[
            {
                "environment": build_env(ctx),
                "image": "plugins/kaniko-ecr",
                "name": "publish",
                "pull": "always",
                "settings": {
                    "access_key": {"from_secret": "ecr_access_key"},
                    "secret_key": {"from_secret": "ecr_secret_key"},
                    "registry": {"from_secret": "ecr_registry"},
                    "repo": ctx.repo.name,
                    "tags": [version_tag(ctx, arch)],
                    "build_args": ["VERSION", "GIT_COMMIT"],
                    "create_repository": True,
                },
            }
        ],
    )


def pipeline_manifest(ctx):
    targets = [version(ctx)]

    # only use "latest" for tagged releases
    if ctx.build.event == "tag":
        targets.append("latest")

    return new_pipeline(
        depends_on=["publish-amd64", "publish-arm64"],
        name="publish-manifest",
        arch="arm64",
        steps=[
            {
                "name": "manifest",
                "image": "public.ecr.aws/kanopy/buildah-plugin:v0.1.1",
                "settings": {
                    "access_key": {"from_secret": "ecr_access_key"},
                    "secret_key": {"from_secret": "ecr_secret_key"},
                    "registry": {"from_secret": "ecr_registry"},
                    "repo": ctx.repo.name,
                    "manifest": {
                        "sources": [
                            version_tag(ctx, "amd64"),
                            version_tag(ctx, "arm64"),
                        ],
                        "targets": targets,
                    },
                },
            },
        ],
    )


def main(ctx):
    pipelines = [pipeline_test(ctx)]

    # only perform image builds for "push" and "tag" events
    if ctx.build.event == "tag" or (
        ctx.build.branch == repo_branch(ctx) and ctx.build.event == "push"
    ):
        pipelines.append(pipeline_build(ctx, "amd64"))
        pipelines.append(pipeline_build(ctx, "arm64"))
        pipelines.append(pipeline_manifest(ctx))

    return pipelines
