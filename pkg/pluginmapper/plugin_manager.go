package pluginmapper

import (
	"errors"
	"github.com/yrn-go/yrn/module/flowmanager"
	"github.com/yrn-go/yrn/pkg/yctx"
)

var (
	_ flowmanager.PluginManager = (*PluginManagerLocal)(nil)
)

type PluginManagerLocal struct{}

func NewPluginManagerLocal() *PluginManagerLocal {
	return &PluginManagerLocal{}
}

func (p *PluginManagerLocal) GetBySlug(ctx *yctx.Context, slug string) (plugin flowmanager.PluginExecutor, err error) {
	var ok bool

	if plugin, ok = mappers[slug]; ok {
		return plugin, nil
	}

	return nil, errors.New("plugin not found")
}
