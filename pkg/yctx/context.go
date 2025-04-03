package yctx

import "context"

type Context struct {
	ctx context.Context
}

func (c *Context) Context() context.Context {
	return c.ctx
}

func NewContext(ctx context.Context) *Context {
	return &Context{
		ctx: ctx,
	}
}
