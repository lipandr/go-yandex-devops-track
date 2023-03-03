package http

import (
	"context"
	"errors"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/lipandr/go-yandex-devops-track/internal/pkg/model"
	"github.com/lipandr/go-yandex-devops-track/internal/server/controller"
)

type ServerHTTP interface {
	Update(w http.ResponseWriter, r *http.Request)
	UpdateJSON(w http.ResponseWriter, r *http.Request)
	GetValue(w http.ResponseWriter, r *http.Request)
	GetJSON(w http.ResponseWriter, r *http.Request)
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

func (h *Handler) UIListAll(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/listAll.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := h.ctl.GetAll(h.ctx)
	// Check client side supports gzip encoding
	if r.Header.Get("Accept-Encoding") == "gzip" {
		w.Header().Set("Content-Encoding", "gzip")
		//gz := gzip.NewWriter(w)
		//defer gz.Close()
		//if err = tmpl.Execute(w, data); err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}
		//if err := gz.Flush(); err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}
	}
	if err = tmpl.Execute(w, data); err != nil {
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
