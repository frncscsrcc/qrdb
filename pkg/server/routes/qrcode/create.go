package qrcode

import (
	"context"
	"encoding/json"
	"net/http"
	"qrdb/pkg/services/qrcode"
)

type CreateQRCode struct{}

type CreateQRCodeInput map[string]string
type CreateQRCodeOutput struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code"`
}

func (route CreateQRCode) ServeHTTP(r *http.Request, ctx context.Context) (int, interface{}) {
	input := make(CreateQRCodeInput)

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return http.StatusBadRequest, map[string]string{"error": "invalid reques"}
	}

	code, err := qrcode.NewService().Create(input)
	errorText := ""
	if err != nil {
		errorText = err.Error()
	}

	output := CreateQRCodeOutput{
		Success: err == nil,
		Error:   errorText,
		Code:    code,
	}

	if err != nil {
		return http.StatusInternalServerError, map[string]string{"error": errorText}
	} else {
		return http.StatusCreated, output
	}

}
