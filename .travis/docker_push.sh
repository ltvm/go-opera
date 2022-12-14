#!/bin/bash

set -euo pipefail

readonly DOCKER_PASSWORD=${DOCKER_PASSWORD:-}

if [[ -z "$DOCKER_PASSWORD" ]]; then
    echo 'DOCKER_PASSWORD is not available, aborting.'
    exit 1
fi

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

docker tag kybernetwork/go-opera:"$TRAVIS_COMMIT" kybernetwork/go-opera:"$TRAVIS_BRANCH"
if [[ -n "$TRAVIS_TAG" ]]; then
    docker tag kybernetwork/go-opera:"$TRAVIS_COMMIT" kybernetwork/go-opera:"$TRAVIS_TAG"
fi

docker push --all-tags kybernetwork/go-opera
