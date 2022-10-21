# build drone-bazelisk-ecr plugin
FROM golang:1.18 AS plugin
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

# setup bazelisk
FROM python:3.8-slim

ENV BAZEL_USER bazel
ENV BAZEL_USER_ID 999
ENV BAZEL_USER_HOME /home/${BAZEL_USER}

ENV BAZELISK_VERSION v1.3.0
ENV BAZELISK_PATH /usr/local/bin/bazel

ENV ECR_LOGIN_VERSION 0.4.0
ENV ECR_LOGIN_PATH /usr/local/bin/docker-credential-ecr-login

RUN groupadd -g ${BAZEL_USER_ID} -r ${BAZEL_USER} \
 && useradd -lmr -u ${BAZEL_USER_ID} -g ${BAZEL_USER} ${BAZEL_USER} \
 && apt-get update && apt-get install -y \
      g++ \
      git \
      unzip \
      wget \
      zip \
 && wget -qO ${BAZELISK_PATH} https://github.com/bazelbuild/bazelisk/releases/download/${BAZELISK_VERSION}/bazelisk-linux-amd64 \
 && chmod +x ${BAZELISK_PATH} \
 && wget -qO ${ECR_LOGIN_PATH} https://amazon-ecr-credential-helper-releases.s3.us-east-2.amazonaws.com/${ECR_LOGIN_VERSION}/linux-amd64/docker-credential-ecr-login \
 && chmod +x ${ECR_LOGIN_PATH}

COPY --from=plugin /go/bin/drone-bazelisk-ecr /usr/local/bin/drone-bazelisk-ecr
COPY --chown=bazel:bazel files/config.json ${BAZEL_USER_HOME}/.docker/config.json

USER ${BAZEL_USER}
ENTRYPOINT ["drone-bazelisk-ecr"]
