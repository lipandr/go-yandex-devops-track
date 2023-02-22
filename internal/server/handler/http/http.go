package http

import (
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/lipandr/go-yandex-devops-track/internal/pkg/model"
	"github.com/lipandr/go-yandex-devops-track/internal/server/controller"
)

type MetHandle interface {
	Update(w http.ResponseWriter, r *http.Request)
	UpdateJSON(w http.ResponseWriter, r *http.Request)
	GetValue(w http.ResponseWriter, r *http.Request)
	GetValueJSON(w http.ResponseWriter, r *http.Request)
	UIListAll(w http.ResponseWriter, r *http.Request)
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

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
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
	name := v[1]
	data, err := getData(v[0], v[2])
	if err != nil {
		http.Error(w, errors.New("bad request").Error(), http.StatusBadRequest)
		return
	}
	if err := h.ctl.Put(h.ctx, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateJSON(w http.ResponseWriter, r *http.Request) {
	var req *model.MetricJSON

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name, data := h.ctl.FromJSON(req)
	if err := h.ctl.Put(h.ctx, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Value updated"))
}

func (h *Handler) GetValue(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	path := r.URL.Path
	trim := strings.TrimPrefix(path, "/value/")
	v := strings.Split(trim, "/")
	if len(v) != 2 {
		http.Error(w, errors.New("bad request").Error(), http.StatusBadRequest)
		return
	}
	name := v[1]
	val, err := h.ctl.Get(h.ctx, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(val))
}

func (h *Handler) GetValueJSON(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)
	w.Header().Set("Content-Type", "application/json")

	var v model.MetricJSON
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	name := v.ID
	val, err := h.ctl.Get(h.ctx, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	switch v.MType {
	case model.TypeGauge:
		tmp, _ := strconv.ParseFloat(val, 64)
		v.Delta = nil
		v.Value = &tmp
	case model.TypeCounter:
		tmp, _ := strconv.ParseInt(val, 10, 64)
		v.Delta = &tmp
		v.Value = nil
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UIListAll(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/listAll.html")
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

func getData(mType string, value string) (*model.Metric, error) {
	var res model.Metric

	res.MType = model.MetricType(mType)
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
