#!/bin/bash
# set -x
set -e

EX=examples/minimal_several_components

make

pushd "$EX" >/dev/null
../../bin/bs clean
../../bin/bs build --build-upstream
popd >/dev/null

dot -Tpng -o /tmp/graphviz.png /tmp/graphviz.dot
feh /tmp/graphviz.png
