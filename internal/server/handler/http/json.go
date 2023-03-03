package http

import (
	"encoding/json"
	"fmt"
	"github.com/lipandr/go-yandex-devops-track/internal/pkg/model"
	"io"
	"net/http"
	"strconv"
)

func (h *Handler) UpdateJSON(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
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
	metric, err := h.GetJSONMetric(req.ID, req.MType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Check client side supports for gzip encoding
	//if r.Header.Get("Accept-Encoding") == "gzip" {
	//	w.Header().Set("Content-Encoding", "gzip")
	//	//gz := gzip.NewWriter(w)
	//	//defer gz.Close()
	//	//err = json.NewEncoder(gz).Encode(val)
	//	//if err != nil {
	//	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	//	return
	//	//}
	//	//return
	//}
	err = json.NewEncoder(w).Encode(metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetJSON(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var metric *model.MetricJSON
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metric, err := h.GetJSONMetric(metric.ID, metric.MType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Check client side supports for gzip encoding
	//if r.Header.Get("Accept-Encoding") == "gzip" {
	//	w.Header().Set("Content-Encoding", "gzip")
	//	//gz := gzip.NewWriter(w)
	//	//defer gz.Close()
	//	//err = json.NewEncoder(gz).Encode(val)
	//	//if err != nil {
	//	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	//	return
	//	//}
	//	//return
	//}
	err = json.NewEncoder(w).Encode(metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetJSONMetric(id string, mType model.MetricType) (*model.MetricJSON, error) {
	res := model.MetricJSON{
		ID:    id,
		MType: mType,
	}
	val, err := h.ctl.Get(h.ctx, id)
	if err != nil {
		return nil, err
	}
	switch mType {
	case model.TypeGauge:
		tmp, _ := strconv.ParseFloat(val, 64)
		res.Value = &tmp
	case model.TypeCounter:
		tmp, _ := strconv.ParseInt(val, 10, 64)
		res.Delta = &tmp
	default:
		return nil, fmt.Errorf("metric type %s not supported", mType)
	}
	return &res, nil
}
