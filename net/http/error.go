package http

import "github.com/anton2920/gofa/context"

func NewError(ctx *context.Context, status Status) bool {
	ctx.NewErrorWithCode(int(status)).S(Status2String[status]).S(": ").S(Status2Reason[status]).Ln().S(ctx.OldError())
	return false
}

func UnauthorizedError(ctx *context.Context) bool {
	return NewError(ctx, StatusUnauthorized)
}

func NotFound(ctx *context.Context) bool {
	return NewError(ctx, StatusNotFound)
}

func ServerError(ctx *context.Context) bool {
	return NewError(ctx, StatusInternalServerError)
}
