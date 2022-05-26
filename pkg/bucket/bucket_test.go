package bucket_test

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/gueckmooh/bs/pkg/bucket"
)

func sl(n int) {
	time.Sleep(3 * time.Second)
}

func sleeper(n int) func() error {
	return func() error {
		sl(n)
		return fmt.Errorf("toto")
	}
}

func TestBucket(t *testing.T) {
	b := bucket.NewBucket(int64(runtime.GOMAXPROCS(0)))
	for i := 0; i < 10; i++ {
		b.Run(sleeper(i))
	}
	err := b.Wait()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	for true {
		if err := b.Error(); err != nil {
		} else {
			break
		}
	}
}
