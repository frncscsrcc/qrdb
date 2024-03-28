package qrcode

import (
	"context"
	"fmt"
	"net/http"
	"qrdb/pkg/server/myctx"

	qrcode "github.com/skip2/go-qrcode"
)

type RenderQR struct{}

type RawPNG []byte

func (route RenderQR) ServeHTTP(r *http.Request, ctx context.Context) (int, interface{}) {
	code := myctx.Stringify(ctx, "code", "")
	page := myctx.Stringify(ctx, "page", "0")
	maxPage := myctx.Stringify(ctx, "max_page", "0")

	var png []byte
	png, err := qrcode.Encode(fmt.Sprintf("%s;%s.%s", code, page, maxPage), qrcode.Medium, 256)

	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, RawPNG(png)
}
