package pluginmapper

import (
	"github.com/yrn-go/yrn/module/flowmanager"
	"github.com/yrn-go/yrn/pkg/pluginhttp"
)

var (
	mappers = map[string]flowmanager.PluginExecutor{
		pluginhttp.SlugHttp: pluginhttp.NewExecutor(),
	}
)
