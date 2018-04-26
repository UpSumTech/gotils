FROM ubuntu:16.04
LABEL "os"="ubuntu" \
  os_version="16.04" \
  golang_version="1.10" \
  author_name="Suman Mukherjee" \
  author_email="sumanmukherjee03@gmail.com" \
  version="$GIT_TAG" \
  build_time="$BUILD_TIME" \
  git_ref="$GIT_REF" \
  build_user="$USER"

ENV BUILD_HOME=/go/src/github.com/sumanmukherjee03/gotils \
  GOPATH="/go" \
  PATH="/go/bin:/usr/lib/go-1.10/bin:$PATH" \
  DEBIAN_FRONTEND=noninteractive \
  CONFIGURE_OPTS="--disable-install-doc" \
  MAKEFLAGS="-j$(($(nproc) + 1))"

RUN apt-key update \
  && apt-get -qq update \
  && chsh -s /bin/bash \
  && rm /bin/sh && ln -sf /bin/bash /bin/sh

RUN set -ex; \
  apt-get install -y --no-install-recommends \
    ca-certificates \
    gcc \
    git \
    autoconf \
    make \
    cmake \
    automake \
    curl \
    wget \
    pkg-config \
    libtool \
    golang-1.10-go; \
  apt-get autoremove -y \
    && apt-get clean -y; \
  mkdir -p $BUILD_HOME

RUN set -ex; \
    [[ ! -d "vendor" ]] \
      && go get -u github.com/golang/dep/cmd/dep \
      && go get -u github.com/mitchellh/gox; \
    go env; \
    command -v dep; \
    command -v gox

WORKDIR $BUILD_HOME
