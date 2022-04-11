#!/bin/bash
set -x
set -e

EX=examples/minimal_lib

make

rm -r "$EX/.build"

(cd $EX && echo -e "----------\n\n\n" && ../../bin/bs && echo -e "\n\n\n----------")

dot -Tpng -o /tmp/graphviz.png /tmp/graphviz.dot
feh /tmp/graphviz.png
