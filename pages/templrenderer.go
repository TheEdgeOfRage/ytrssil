package pages

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
)

type TemplRenderer struct {
	Ctx       context.Context
	Component templ.Component
}

func (t TemplRenderer) Render(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if t.Component != nil {
		return t.Component.Render(t.Ctx, w)
	}
	return nil
}

func (t TemplRenderer) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}
