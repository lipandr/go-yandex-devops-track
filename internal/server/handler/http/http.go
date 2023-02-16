package http

import (
	"context"
	"encoding/json"
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
	PutMetricJSON(w http.ResponseWriter, r *http.Request)
	GetMetricValue(w http.ResponseWriter, r *http.Request)
	ListAllMetrics(w http.ResponseWriter, r *http.Request)
	GetMetricValueJSON(w http.ResponseWriter, r *http.Request)
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
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

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
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) PutMetricJSON(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	var v *model.MetricJSON

	buf := io.NopCloser(r.Body)
	if err := json.NewDecoder(buf).Decode(&v); err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req, err := convertFromJSON(*v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.ctl.Put(h.ctx, req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetMetricValue(w http.ResponseWriter, r *http.Request) {
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
	req, _ := getMetric(v[0], v[1], "")

	val, err := h.ctl.Get(h.ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(val))
}
func (h *Handler) GetMetricValueJSON(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	var v model.MetricJSON
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	req := model.Metric{
		ID:    v.ID,
		MType: v.MType,
	}
	val, err := h.ctl.Get(h.ctx, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	switch v.MType {
	case model.TypeGauge:
		tmp, _ := strconv.ParseFloat(val, 64)
		v.Value = &tmp
	case model.TypeCounter:
		tmp, _ := strconv.ParseInt(val, 10, 64)
		v.Delta = &tmp
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ListAllMetrics(w http.ResponseWriter, r *http.Request) {
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
func convertFromJSON(data model.MetricJSON) (*model.Metric, error) {
	if data.Delta != nil {
		return getMetric(string(data.MType), data.ID, strconv.Itoa(int(*data.Delta)))
	}
	if data.Value != nil {
		return getMetric(string(data.MType), data.ID, strconv.FormatFloat(*data.Value, 'f', -1, 64))
	}
	return nil, errors.New("bad request")
}

func checkType(mType string) bool {
	switch mType {
	case model.TypeGauge, model.TypeCounter:
		return true
	default:
		return false
	}
}
