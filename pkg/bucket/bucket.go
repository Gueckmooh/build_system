package bucket

import (
	"container/list"
	"context"
	"sync"

	"github.com/gueckmooh/bs/pkg/bucket/semaphore"
)

type Bucket struct {
	maxWorkers int64
	sema       *semaphore.Weighted
	mutex      sync.Mutex
	ctx        context.Context
	errors     list.List
}

func NewBucket(nb int64) *Bucket {
	return &Bucket{
		maxWorkers: nb,
		sema:       semaphore.NewWeighted(nb),
		ctx:        context.TODO(),
	}
}

func (b *Bucket) Run(f func() error) error {
	if err := b.sema.Acquire(b.ctx, 1); err != nil {
		return err
	}
	go func() {
		defer b.sema.Release(1)
		err := f()
		if err != nil {
			b.mutex.Lock()
			b.errors.PushBack(err)
			b.mutex.Unlock()
		}
	}()
	return nil
}

func (b *Bucket) RunFailIfError(f func() error) error {
	if err := b.sema.Acquire(b.ctx, 1); err != nil {
		return err
	}
	if err := b.Error(); err != nil {
		defer b.sema.Release(1)
		return err
	}
	go func() {
		defer b.sema.Release(1)
		err := f()
		if err != nil {
			b.mutex.Lock()
			b.errors.PushBack(err)
			b.mutex.Unlock()
		}
	}()
	return nil
}

func (b *Bucket) Error() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.errors.Len() > 0 {
		err := b.errors.Front()
		b.errors.Remove(err)
		return err.Value.(error)
	}
	return nil
}

func (b *Bucket) Wait() error {
	return b.sema.Wait(b.ctx)
}
