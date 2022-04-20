module github.com/gueckmooh/bs

go 1.18

require (
	github.com/alessio/shellescape v1.4.1
	github.com/yuin/gopher-lua v0.0.0-20210529063254-f4c35e4016d9
)

replace github.com/yuin/gopher-lua => ./internal/gopher-lua
