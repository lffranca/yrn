package yctx

import "context"

type Context struct {
	ctx context.Context
}

func (c *Context) Context() context.Context {
	return c.ctx
}
