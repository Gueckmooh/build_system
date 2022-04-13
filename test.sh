#!/bin/bash
# set -x
set -e

PATH=$PATH:$(realpath bin)
EX=examples/basic

make

pushd "$EX" >/dev/null
bs clean
popd >/dev/null

bs build --build-upstream -C "$EX" --verbose
"$EX"/.build/bin/hello_exe
