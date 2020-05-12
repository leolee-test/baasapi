#!/usr/bin/env sh

binary="baasapi-$1-$2"

mkdir -p dist

docker run --rm -tv "$(pwd)/api:/src" -e BUILD_GOOS="$1" -e BUILD_GOARCH="$2" baasapi/golang-builder:cross-platform /src/cmd/baasapi

mv "api/cmd/baasapi/$binary" dist/
#sha256sum "dist/$binary" > baasapi-checksum.txt
