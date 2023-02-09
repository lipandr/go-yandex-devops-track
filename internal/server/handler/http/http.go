package http

import (
	"context"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/lipandr/go-yandex-devops-track/internal/server/controller"
)

type MetHandle interface {
	PutMetric(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	ctx context.Context
	ctl *controller.Controller
}

func New(ctx context.Context, controller *controller.Controller) *Handler {
	return &Handler{
		ctx: ctx,
		ctl: controller,
	}
}

func (h *Handler) PutMetric(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	trim := strings.TrimPrefix(path, "/update/")
	v := strings.Split(trim, "/")
	if err := h.ctl.Put(h.ctx, v[1], v[0], v[2]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)
}

func (h *Handler) GetMetricValue(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	trim := strings.TrimPrefix(path, "/value/")
	v := strings.Split(trim, "/")
	val, err := h.ctl.Get(h.ctx, v[0], v[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(val))
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)
}

func (h *Handler) ListAllMetrics(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("listAllMetrics.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := h.ctl.GetAll(h.ctx)
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)
}
