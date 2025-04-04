package flowmanager

import (
	"github.com/yrn-go/yrn/pkg/yctx"
	"log/slog"
)

type (
	PluginExecutor interface {
		Do(ctx *yctx.Context, schemaInputs string, previousPluginResponse any, responseSharedForAll map[string]any) (output any, err error)
	}
	PluginManager interface {
		GetBySlug(ctx *yctx.Context, slug string) (plugin PluginExecutor, err error)
	}
	FlowExecutor struct {
		flowReaderRepository FlowReaderRepository
		pluginManager        PluginManager
	}
)

func NewFlowExecutor(
	flowReaderRepository FlowReaderRepository,
	pluginManager PluginManager,
) *FlowExecutor {
	return &FlowExecutor{
		flowReaderRepository,
		pluginManager,
	}
}

func (f *FlowExecutor) Do(ctx *yctx.Context, flowId string, eventRequestData any) (response any, err error) {

	var (
		flow                 *Flow
		responseSharedForAll = make(map[string]any)
	)

	flow, err = f.flowReaderRepository.GetById(ctx, flowId)
	if err != nil {
		return
	}

	for index, pluginInfo := range flow.Plugins {
		var (
			pluginExecutor PluginExecutor
			pluginData     = response
			pluginResponse any
		)

		if index == 0 {
			pluginData = eventRequestData
		}

		pluginExecutor, err = f.pluginManager.GetBySlug(ctx, pluginInfo.Slug)
		if err != nil {
			slog.Error(
				"flow plugin get by slug",
				slog.Any("plugin_info", pluginInfo),
				slog.Any("error", err),
			)

			if !pluginInfo.ContinueEvenWithError {
				return
			}

			continue
		}

		pluginResponse, err = pluginExecutor.Do(ctx, pluginInfo.SchemaInput, pluginData, responseSharedForAll)
		if err != nil {
			slog.Error(
				"flow plugin error execute",
				slog.Any("plugin_info", pluginInfo),
				slog.Any("error", err),
			)

			if !pluginInfo.ContinueEvenWithError {
				return
			}

			continue
		}

		if pluginInfo.ShareResponseWithAllPlugins {
			responseSharedForAll[pluginInfo.Id] = pluginResponse
		}

		response = pluginResponse
	}

	return
}
