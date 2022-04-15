package bucket_test

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/gueckmooh/bs/pkg/bucket"
)

func sl(n int) {
	fmt.Printf("sleep %d\n", n)
	time.Sleep(3 * time.Second)
	fmt.Printf("end sleep %d\n", n)
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
	fmt.Println("Waiting")
	err := b.Wait()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Sleeping a bit more")
	time.Sleep(time.Second)
	fmt.Printf("End\n")
	for true {
		if err := b.Error(); err != nil {
			fmt.Printf("-> %s\n", err.Error())
		} else {
			break
		}
	}
}
