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

type ConnectorInput[T any] struct {
	Body T `json:"body"`
}

type ConnectorOutput[T any] struct {
	ProcessingStatus ProcessingStatus `json:"processingStatus"`
	Body             *T               `json:"body"`
}

type Connector[Input any, Output any] interface {
	Schema(ctx context.Context) map[string]any
	Do(ctx context.Context, input *ConnectorInput[Input]) (output *ConnectorOutput[Output], err error)
}
