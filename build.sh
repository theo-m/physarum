#!/usr/bin/env bash

set -eEx
trap 'catch' ERR EXIT
catch() {
  rm Dockerfile .dockerignore physarum-server
}

go generate

pushd js
./node_modules/webpack-cli/bin/cli.js && mv bundle.js ../public
popd

# https://github.com/hashicorp/vault/issues/3417
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o physarum-server

cat << EOF > Dockerfile
FROM alpine
RUN apk add --no-cache ffmpeg
COPY physarum-server .
COPY physarum.gohtml .
COPY physarum.proto .
ADD public public/
CMD ./physarum-server
EOF

cat << EOF > .dockerignore
**/node_modules
out/
EOF

docker build .
