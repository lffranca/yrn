package ybase

import "context"

type (
	FlowExecuteDataInput struct{}

	FlowExecuteOutput struct{}

	ConnectorProcessingDataRepository interface {
		SavingProcessingData(ctx context.Context, request, response any) (err error)
	}

	FlowExecutor struct {
		connectorProcessingData ConnectorProcessingDataRepository
	}
)

func (executor *FlowExecutor) Execute(ctx context.Context, flow *Flow, inputData *FlowExecuteDataInput) (output *FlowExecuteOutput, err error) {
	return
}
