ARG NON_ROOT_UID=1001
ARG NON_ROOT_GID=1001
ARG NON_ROOT_USER=default
FROM sumanmukherjee03/golang:onbuild-1.10.0
LABEL git_tag="$GIT_TAG" \
  build_time="$BUILD_TIME" \
  git_ref="$GIT_REF" \
  github_user="$GITHUB_USERNAME" \
  build_user="$BUILD_USER" \
  repo_name="$REPO_NAME"

RUN set -ex; \
  ls -lah /var/data/build

CMD []
