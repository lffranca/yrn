package ybase

import (
	"context"
)

type ProcessingStatus string

const (
	ProcessingStatusPending    ProcessingStatus = "PENDING"
	ProcessingStatusInProgress ProcessingStatus = "IN_PROGRESS"
	ProcessingStatusCompleted  ProcessingStatus = "COMPLETED"
	ProcessingStatusFailed     ProcessingStatus = "FAILED"
)

type ConnectorInput struct {
	Body any `json:"body"`
}

type ConnectorOutput struct {
	ProcessingStatus ProcessingStatus `json:"processingStatus"`
	Body             any              `json:"body"`
}

type Connector interface {
	Schema(ctx context.Context) map[string]any
	Do(ctx context.Context, input *ConnectorInput) (output *ConnectorOutput, err error)
}
