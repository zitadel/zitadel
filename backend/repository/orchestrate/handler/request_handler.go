package handler

import "context"

type Handle[Req, Res any] func(ctx context.Context, request Req) (res Res, err error)

type Decorate[Req, Res any] func(ctx context.Context, request Req, handle Handle[Req, Res]) (res Res, err error)

func NewChained[Req, Res any](handle Handle[Req, Res], next Handle[Res, Res]) Handle[Req, Res] {
	return func(ctx context.Context, request Req) (res Res, err error) {
		res, err = handle(ctx, request)
		if err != nil {
			return res, err
		}
		if next == nil {
			return res, nil
		}
		return next(ctx, res)
	}
}

func NewDecorated[Req, Res any](decorate Decorate[Req, Res], handle Handle[Req, Res]) Handle[Req, Res] {
	return func(ctx context.Context, request Req) (res Res, err error) {
		return decorate(ctx, request, handle)
	}
}
