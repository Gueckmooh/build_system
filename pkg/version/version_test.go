package version_test

import (
	"testing"

	"github.com/gueckmooh/bs/pkg/version"
)

func TestVersion1(t *testing.T) {
	var err error
	_, err = version.ParseVersionHash("v0.1.0-1-ga12f880")
	if err != nil {
		t.Fatal(err)
	}
	_, err = version.ParseVersionHash("v0.1.0")
	if err != nil {
		t.Fatal(err)
	}
	_, err = version.ParseVersionHash("v0.0.0")
	if err != nil {
		t.Fatal(err)
	}
	_, err = version.ParseVersionHash("v0.1-1-ga12f880")
	if err == nil {
		t.Fatal("An err should have been returned")
	}
}
