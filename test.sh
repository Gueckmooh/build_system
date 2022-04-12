#!/bin/bash
# set -x
set -e

PATH=$PATH:$(realpath bin)
EX=examples/basic/sources/hello

make

pushd "$EX" >/dev/null
bs clean
bs build --build-upstream
popd >/dev/null
