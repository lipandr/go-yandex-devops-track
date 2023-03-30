package http

import (
	"context"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/lipandr/go-yandex-devops-track/internal/pkg/model"
	"github.com/lipandr/go-yandex-devops-track/internal/server/controller"
)

type MetHandle interface {
	PutMetric(w http.ResponseWriter, r *http.Request)
	GetMetricValue(w http.ResponseWriter, r *http.Request)
	ListAllMetrics(w http.ResponseWriter, r *http.Request)
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
	if !checkType(v[0]) {
		http.Error(w, errors.New("not implemented").Error(), http.StatusNotImplemented)
		return
	}
	if len(v) < 3 {
		http.Error(w, errors.New("not found").Error(), http.StatusNotFound)
		return
	}
	req, err := getMetric(v[0], v[1], v[2])
	if err != nil {
		http.Error(w, errors.New("bad request").Error(), http.StatusBadRequest)
		return
	}
	if err := h.ctl.Put(h.ctx, req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)
}

func (h *Handler) GetMetricValue(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	trim := strings.TrimPrefix(path, "/value/")
	v := strings.Split(trim, "/")
	if len(v) != 2 {
		http.Error(w, errors.New("bad request").Error(), http.StatusBadRequest)
		return
	}
	req, _ := getMetric(v[0], v[1], "")

	val, err := h.ctl.Get(h.ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(val))
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)
}

func (h *Handler) ListAllMetrics(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./app/listAllMetrics.html")
	if err != nil {
		log.Printf("error parsing template: %v", err)
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

func getMetric(mType string, id string, value string) (*model.Metric, error) {
	var res model.Metric
	res.MType = model.MetricType(mType)
	res.ID = id
	if value != "" {
		switch mType {
		case model.TypeGauge:
			fl, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			res.Value = fl
		case model.TypeCounter:
			in, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
			res.Delta = int64(in)
		}
	}
	return &res, nil
}

func checkType(mType string) bool {
	switch mType {
	case model.TypeGauge, model.TypeCounter:
		return true
	default:
		return false
	}
}
