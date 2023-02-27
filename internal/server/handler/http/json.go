package http

import (
	"encoding/json"
	"github.com/lipandr/go-yandex-devops-track/internal/pkg/model"
	"io"
	"net/http"
	"strconv"
)

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
		http.Error(w, err.Error(), http.StatusOK)
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
