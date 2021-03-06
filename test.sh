#!/bin/bash
# set -x
set -e

PATH=$PATH:$(realpath bin)
EX=examples/basic_with_extra_sources

make

pushd "$EX" >/dev/null
bs clean
popd >/dev/null

bs build --build-upstream -C "$EX" --verbose -P Linux

echo -e "\nRunning..."
"$EX"/.build/bin/hello_exe
