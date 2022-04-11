#!/bin/bash
# set -x
set -e

EX=examples/minimal_lib

make

rm -rf "$EX/.build"

pushd "$EX" >/dev/null
../../bin/bs build
popd >/dev/null

# dot -Tpng -o /tmp/graphviz.png /tmp/graphviz.dot
# feh /tmp/graphviz.png
