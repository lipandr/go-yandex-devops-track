package http

import (
	"context"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lipandr/go-yandex-devops-track/internal/pkg/model"
	"github.com/lipandr/go-yandex-devops-track/internal/server/controller"
	"github.com/lipandr/go-yandex-devops-track/internal/server/storage/memory"
)

func TestHandler_GetValue(t *testing.T) {
	repo := memory.New()
	_ = repo.Put(context.Background(), "test", &model.Metric{
		MType: model.TypeCounter,
		Delta: 123,
	})
	cnt := controller.NewMemoryRepo(repo)
	handler := New(context.Background(), cnt)
	type want struct {
		request     string
		statusCode  int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "get metric value",
			want: want{
				request:     "/value/counter/test",
				statusCode:  http.StatusOK,
				response:    "123",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "get metric value not found",
			want: want{
				request:     "/value/counter/NotFound",
				statusCode:  http.StatusNotFound,
				response:    "not found\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "get metric not implemented error",
			want: want{
				request:     "/value/unknown/PollCount",
				statusCode:  http.StatusNotFound,
				response:    "not found\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "get metric bad request error",
			want: want{
				request:     "/value/gauge",
				statusCode:  http.StatusBadRequest,
				response:    "bad request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.want.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.GetValue)
			h.ServeHTTP(w, req)

			require.Equalf(t, tt.want.statusCode, w.Code,
				"handler.GetValue() = %v, want %v", w.Code, tt.want.statusCode)
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(req.Body)

			require.Equalf(t, tt.want.response, w.Body.String(),
				"handler.GetValue() = %v, want %v", w.Body.String(), tt.want.response)
			require.Equalf(t, tt.want.contentType, w.Header().Get("Content-Type"),
				"handler.GetValue() = %v, want %v", w.Header().Get("Content-Type"), tt.want.contentType)
		})
	}
}

func TestHandler_Update(t *testing.T) {
	repo := memory.New()
	cnt := controller.NewMemoryRepo(repo)
	handler := New(context.Background(), cnt)
	type want struct {
		request     string
		statusCode  int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "put metric",
			want: want{
				request:     "/update/counter/PollCount/123",
				statusCode:  http.StatusOK,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "put metric bad request error",
			want: want{
				request:     "/update/gauge/test/none",
				statusCode:  http.StatusBadRequest,
				response:    "bad request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "put metric not implemented error",
			want: want{
				request:     "/update/unknown/PollCount/123",
				statusCode:  http.StatusNotImplemented,
				response:    "not implemented\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.want.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.Update)
			h.ServeHTTP(w, req)

			require.Equalf(t, tt.want.statusCode, w.Code, "handler.Update() = %v, want %v", w.Code, tt.want.statusCode)
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(req.Body)

			require.Equalf(t, tt.want.response, w.Body.String(),
				"handler.Update() = %v, want %v", w.Body.String(), tt.want.response)
			require.Equalf(t, tt.want.contentType, w.Header().Get("Content-Type"),
				"handler.Update() = %v, want %v", w.Header().Get("Content-Type"), tt.want.contentType)
		})
	}
}
