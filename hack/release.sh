#!/bin/sh

COMMIT_SHA=$(git rev-parse --verify --quiet HEAD)
GIT_BRANCH=$(git rev-parse --symbolic-full-name --verify --quiet --abbrev-ref HEAD)
GIT_REPO="rkrmr33/onka"
VERSION=$(cat ./VERSION)

if [[ -z "$PRERELEASE" ]]; then
    PRERELEASE=false
fi


if [[ "$(git branch -r --contains $COMMIT_SHA)" != "" ]]; then
    echo "local branch is up to date with remote branch"
else
    echo "ERROR: local brach is not up to date with remote branch"
    exit 1
fi

echo "on release branch: $GIT_BRANCH"
echo "running: gh release create $VERSION --repo $GIT_REPO --notes $VERSION -t $VERSION --target $GIT_BRANCH --prerelease=$PRERELEASE"

if [[ "$DRY_RUN" == "1" ]]; then
    exit 0
fi

gh release create $VERSION --repo $GIT_REPO --notes $VERSION -t $VERSION --target $GIT_BRANCH --prerelease=$PRERELEASE
