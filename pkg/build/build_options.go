package build

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
