package gcc_test

import (
	"testing"

	"github.com/gueckmooh/bs/pkg/compiler/gcc"
)

func TestParseM(t *testing.T) {
	toParse := `hello.o: \
 hello.cpp \
 hello.hpp`
	target, sources, err := gcc.ParseMOutput(toParse)
	if err != nil {
		t.Fail()
	}
	if target != "hello.o" {
		t.Fail()
	}
	if sources[0] != "hello.cpp" || sources[1] != "hello.hpp" {
		t.Fail()
	}
}
