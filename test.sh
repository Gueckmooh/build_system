#!/bin/bash
set -x
set -e

EX=examples/minimal

make

rm -rf "$EX/.build"

(cd $EX && echo -e "----------\n\n\n" && ../../bin/bs build && echo -e "\n\n\n----------")

dot -Tpng -o /tmp/graphviz.png /tmp/graphviz.dot
feh /tmp/graphviz.png
