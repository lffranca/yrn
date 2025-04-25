package flowmanager

import (
	"github.com/yrn-go/yrn/pkg/yctx"
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
		flow         *Flow
		eventManager = NewEventManager(f.pluginManager)
	)

	flow, err = f.flowReaderRepository.GetById(ctx, flowId)
	if err != nil {
		return
	}

	for _, pluginInfo := range flow.Plugins {
		if err = eventManager.Register(pluginInfo); err != nil {
			return
		}
	}

	return eventManager.Execute(ctx, flow.FirstPluginToRun, eventRequestData)
}
