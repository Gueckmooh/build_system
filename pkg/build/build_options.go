package build

import "github.com/gueckmooh/bs/pkg/lua"

type BuildOption func(b *Builder)

// func WithBuildUpstream(b *Builder) {
// 	b.buildUpstream = true
// }

func WithAlwaysBuild(b *Builder) {
	b.alwaysBuild = true
}

func WithProfile(s string) BuildOption {
	return func(b *Builder) {
		b.profile = s
	}
}

func WithPlatform(s string) BuildOption {
	return func(b *Builder) {
		b.platform = s
	}
}

func WithJobs(j int) BuildOption {
	return func(b *Builder) {
		b.jobs = j
	}
}

func WithLuaContect(C *lua.LuaContext) BuildOption {
	return func(b *Builder) {
		b.C = C
	}
}
