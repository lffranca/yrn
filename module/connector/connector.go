package connector

import "context"

type Type struct {
	Id   string
	Name string
}

type Connector interface {
	Schema()
	Process(ctx context.Context, data []byte) ([]byte, error)
}
