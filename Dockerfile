FROM ubuntu:16.04
LABEL "os"="ubuntu" \
  os_version="16.04" \
  golang_version="1.10" \
  author_name="Suman Mukherjee" \
  author_email="sumanmukherjee03@gmail.com" \
  version="$GIT_TAG" \
  build_time="$BUILD_TIME" \
  git_ref="$GIT_REF" \
  github_user="$GITHUB_USERNAME" \
  build_user="$USER"

ENV BUILD_HOME="/go/src/github.com/$GITHUB_USERNAME/gotils" \
  BUILD_DATA="/var/data" \
  GOPATH="/go" \
  PATH="/go/bin:/usr/lib/go-1.10/bin:$PATH" \
  DEBIAN_FRONTEND="noninteractive" \
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
    tar \
    zip \
    unzip \
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
    && apt-get clean -y

RUN set -ex; \
    mkdir -p $BUILD_HOME; \
    [[ ! -d "vendor" ]] \
      && go get -u github.com/golang/dep/cmd/dep \
      && go get -u github.com/mitchellh/gox; \
    go env; \
    command -v dep; \
    command -v gox

WORKDIR $BUILD_HOME

COPY Gopkg.toml Gopkg.toml
COPY Gopkg.lock Gopkg.lock
COPY main.go main.go
COPY cmd cmd/

RUN set -ex; \
  dep ensure; \
  CGO_ENABLED=0 gox -osarch='linux/amd64 linux/386 darwin/amd64 darwin/386' -rebuild -tags='netgo' -ldflags='-w -extldflags "-static"'; \
  mv gotils_linux_* $BUILD_DATA; \
  mv gotils_darwin_* $BUILD_DATA; \
  mkdir -p $BUILD_DATA; \
  cd $BUILD_DATA; \
  tar czf gotils.tar.gz gotils_linux_* gotils_darwin_*; \
  rm gotils_linux_* gotils_darwin_*
