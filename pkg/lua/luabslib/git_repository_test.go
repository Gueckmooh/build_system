package luabslib_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gueckmooh/bs/pkg/lua/luabslib"
	lua "github.com/yuin/gopher-lua"
)

func TestGitRepositoryClone(t *testing.T) {
	tmpdir := t.TempDir()
	err := os.Chdir(tmpdir)
	if err != nil {
		t.Fatal(err)
	}
	L := lua.NewState()
	defer L.Close()
	luabslib.RegisterTypes(L)
	if err := L.DoString(`
gr = GitRepository.new{ url = "https://github.com/gueckmooh/build_system",
                        path = "here" }
gr:Clone()
`); err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	stat, err := os.Stat("here")
	if err != nil {
		t.Fatal(err)
	}
	if !stat.IsDir() {
		t.Fatal("here must be a directory")
	}
}

func TestGitRepositoryCloneRevision(t *testing.T) {
	tmpdir := t.TempDir()
	err := os.Chdir(tmpdir)
	if err != nil {
		t.Fatal(err)
	}
	L := lua.NewState()
	defer L.Close()
	luabslib.RegisterTypes(L)
	if err := L.DoString(`
gr = GitRepository.new{ url = "https://github.com/gueckmooh/build_system",
                        path = "here",
                        revision = "e4c8f6474dcebf9188fb0395f7f870e4807dfa91" }
gr:Clone()
gr:Checkout()
`); err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	stat, err := os.Stat("here")
	if err != nil {
		t.Fatal(err)
	}
	if !stat.IsDir() {
		t.Fatal("here must be a directory")
	}
}
