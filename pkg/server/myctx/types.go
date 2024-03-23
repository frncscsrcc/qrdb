package myctx

import "context"

type param string
type roles string

func Param(s string) param {
	return param(s)
}

func Roles(s string) roles {
	return roles(s)
}

func Stringify(ctx context.Context, code string, _default string) string {
	x := ctx.Value(param(code))
	if str, ok := x.(string); ok {
		return str
	} else {
		return _default
	}
}
