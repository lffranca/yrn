package flowmanager

import (
	"github.com/yrn-go/yrn/pkg/yctx"
	"log/slog"
)

type (
	JSONSchemaValidator interface {
		Validate(ctx *yctx.Context, schema string, body []byte) (err error)
	}
	PluginExecutor interface {
		Execute(ctx *yctx.Context, plugin *Plugin, data []byte, responseSharedForAll map[string][]byte) (response []byte, err error)
	}
	FlowExecutor struct {
		jsonSchemaValidator  JSONSchemaValidator
		flowReaderRepository FlowReaderRepository
		pluginExecutor       PluginExecutor
	}
)

func (f *FlowExecutor) Do(ctx *yctx.Context, flowId string, data []byte) (err error) {

	var (
		flow                   *Flow
		previousPluginResponse []byte
		responseSharedForAll   = make(map[string][]byte)
	)

	flow, err = f.flowReaderRepository.GetById(ctx, flowId)
	if err != nil {
		return
	}

	for index, plugin := range flow.Plugins {
		var (
			pluginData     = previousPluginResponse
			pluginResponse []byte
		)

		if index == 0 {
			pluginData = data
		}

		if err = f.jsonSchemaValidator.Validate(ctx, plugin.Schema, pluginData); err != nil {
			slog.Error(
				"flow plugin error validate",
				slog.Any("plugin", plugin),
				slog.Any("error", err),
			)

			if !plugin.ContinueEvenWithError {
				return
			}
		}

		pluginResponse, err = f.pluginExecutor.Execute(ctx, &plugin, pluginData, responseSharedForAll)
		if err != nil {
			slog.Error(
				"flow plugin error execute",
				slog.Any("plugin", plugin),
				slog.Any("error", err),
			)

			if !plugin.ContinueEvenWithError {
				return
			}
		}

		if plugin.ShareResponseWithAllPlugins {
			responseSharedForAll[plugin.Slug] = pluginResponse
		}

		previousPluginResponse = pluginResponse
	}

	return
}
