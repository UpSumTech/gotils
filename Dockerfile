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
