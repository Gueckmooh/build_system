package compiler

import (
	"github.com/gueckmooh/bs/pkg/bucket"
)

type Scheduler struct {
	compiler Compiler
	njobs    int64
	b        *bucket.Bucket
}

func NewScheduler(c Compiler, j int64) *Scheduler {
	s := &Scheduler{
		compiler: c,
		njobs:    j,
		b:        nil,
	}
	if j > 1 {
		s.b = bucket.NewBucket(j)
	}
	return s
}

func (s *Scheduler) CompileFile(target, source string) error {
	if s.njobs > 1 {
		return s.b.RunFailIfError(func() error {
			return s.compiler.CompileFile(target, source)
		})
	} else {
		return s.compiler.CompileFile(target, source)
	}
}

func (s *Scheduler) LinkFiles(target string, sources ...string) error {
	if s.njobs > 1 {
		if err := s.b.Wait(); err != nil {
			return err
		}
		if err := s.b.Error(); err != nil {
			return err
		}
	}
	return s.compiler.LinkFiles(target, sources...)
}
