#!/usr/bin/env bash

set -eEx

protoc --python_out=physarum --mypy_out=physarum --js_out=import_style=commonjs,binary:js physarum.proto
pushd js
./node_modules/webpack-cli/bin/cli.js
popd

python -m physarum.server