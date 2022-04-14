#!/bin/bash
# set -x
set -e

PATH=$PATH:$(realpath bin)
EX=examples/basic

make

pushd "$EX" >/dev/null
bs clean
popd >/dev/null

bs build --build-upstream -C "$EX" --verbose -P Debug

echo -e "\nRunning..."
"$EX"/.build/bin/hello_exe
