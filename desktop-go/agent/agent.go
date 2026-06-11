package agent

import (
	"context"
	"embed"
	"fmt"
	"yekonga-builder/console"
	"yekonga-builder/types"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:templates
var assets embed.FS

type Agent struct {
	ApiKey string
	ctx    *context.Context
}

func NewAgent(apiKey string, ctx *context.Context) *Agent {
	return &Agent{
		ApiKey: apiKey,
		ctx:    ctx,
	}
}

func (a *Agent) Emit(p types.ScaffoldProgress) {
	runtime.EventsEmit(*a.ctx, "scaffold:progress", p)
}

func (a *Agent) getSystemInstruction(name string) string {
	value, err := assets.ReadFile(fmt.Sprintf("template/%s.%s", name, "md"))

	if err != nil {
		console.Log(err.Error())
	}

	return string(value)
}
