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
	res := fmt.Sprintf("Metric %s updated", name)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	//// Check client side supports gzip encoding
	//if r.Header.Get("Accept-Encoding") == "gzip" {
	//	w.Header().Set("Content-Encoding", "gzip")
	//	gz := gzip.NewWriter(w)
	//	defer gz.Close()
	//	if _, err := gz.Write([]byte(res)); err != nil {
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//		return
	//	}
	//	if err := gz.Flush(); err != nil {
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//		return
	//	}
	//}
	_, _ = w.Write([]byte(res))
}

func (h *Handler) GetValueJSON(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var metric model.MetricJSON
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	val, err := h.ctl.Get(h.ctx, metric.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	switch metric.MType {
	case model.TypeGauge:
		tmp, _ := strconv.ParseFloat(val, 64)
		metric.Value = &tmp
	case model.TypeCounter:
		tmp, _ := strconv.ParseInt(val, 10, 64)
		metric.Delta = &tmp
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//// Check client side supports for gzip encoding
	//if r.Header.Get("Accept-Encoding") == "gzip" {
	//	w.Header().Set("Content-Encoding", "gzip")
	//	gz := gzip.NewWriter(w)
	//	defer gz.Close()
	//	err = json.NewEncoder(gz).Encode(val)
	//	if err != nil {
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//		return
	//	}
	//	return
	//}
	err = json.NewEncoder(w).Encode(metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
