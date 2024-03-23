package routes

import (
	"context"
	"net/http"
)

type Status struct{}

func (route Status) AllowsRoles() []string { return []string{} }

func (route Status) AllowsMethods() []string { return []string{} }

func (route Status) ServeHTTP(r *http.Request, ctx context.Context) (int, interface{}) {
	return http.StatusOK, map[string]string{"status": "ok"}
}
