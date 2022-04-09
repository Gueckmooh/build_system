package build

type BuildOption func(b *Builder)

func WithBuildUpstream(b *Builder) {
	b.buildUpstream = true
}
