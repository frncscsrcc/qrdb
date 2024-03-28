package qrcode

import (
	"context"
	"net/http"
	"qrdb/pkg/server/myctx"
	"qrdb/pkg/services/qrcode"
)

type GetData struct{}

type GetDataOutput struct {
	Success bool              `json:"success"`
	Error   string            `json:"error"`
	Data    map[string]string `json:"data"`
}

func (route GetData) ServeHTTP(r *http.Request, ctx context.Context) (int, interface{}) {
	code := myctx.Stringify(ctx, "code", "")

	if code == "" {
		return http.StatusBadRequest, "missing code"
	}

	data, err := qrcode.NewService().Get(code)
	if len(data) == 0 {
		return http.StatusNotFound, "not found"
	}

	errorText := ""
	if err != nil {
		errorText = err.Error()
	}

	output := GetDataOutput{
		Success: err == nil,
		Error:   errorText,
		Data:    data,
	}

	if err != nil {
		return http.StatusInternalServerError, errorText
	} else {
		return http.StatusOK, output
	}

}
