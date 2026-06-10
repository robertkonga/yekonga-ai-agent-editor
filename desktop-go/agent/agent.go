package agent

import "embed"

//go:embed all:templates
var assets embed.FS

type Agent struct {
}
