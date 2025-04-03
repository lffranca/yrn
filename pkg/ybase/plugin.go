package ybase

import (
	"context"
)

type ProcessingStatus string

const (
	ProcessingStatusSucceeded ProcessingStatus = "SUCCEEDED"
	ProcessingStatusFailed    ProcessingStatus = "FAILED"
)

type PluginInput struct {
	Body []byte `json:"body"`
}

type PluginOutput struct {
	ProcessingStatus ProcessingStatus `json:"processingStatus"`
	Body             []byte           `json:"body"`
}

type Plugin interface {
	Schema(ctx context.Context) map[string]any
	Do(ctx context.Context, input *PluginInput) (output *PluginOutput, err error)
}
