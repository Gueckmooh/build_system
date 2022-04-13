package build

type BuildOption func(b *Builder)

// func WithBuildUpstream(b *Builder) {
// 	b.buildUpstream = true
// }

func WithAlwaysBuild(b *Builder) {
	b.alwaysBuild = true
}
