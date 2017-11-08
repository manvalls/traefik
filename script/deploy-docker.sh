#!/usr/bin/env bash
set -e

if [ -n "$TRAVIS_COMMIT" ]; then
  echo "Deploying PR..."
else
  echo "Skipping deploy PR"
  exit 0
fi

# create docker image manvalls/traefik
echo "Updating docker manvalls/traefik image..."
docker login -u $DOCKER_USER -p $DOCKER_PASS
docker tag manvalls/traefik manvalls/traefik:${TRAVIS_COMMIT}
docker push manvalls/traefik:${TRAVIS_COMMIT}
docker tag manvalls/traefik manvalls/traefik:experimental
docker push manvalls/traefik:experimental

echo "Deployed"
