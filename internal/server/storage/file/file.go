package file

import (
	"encoding/json"
	"os"

	"github.com/lipandr/go-yandex-devops-track/internal/pkg/model"
)

type repository interface {
	Write(metric *model.MetricJSON) error
	Read() (*model.MetricJSON, error)
	Close() error
}

type Writer struct {
	repository
	file    *os.File
	encoder *json.Encoder
}

func NewFileWriter(filename string) (*Writer, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	return &Writer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (w *Writer) Write(metric *model.MetricJSON) error {
	return w.encoder.Encode(&metric)
}

func (w *Writer) Close() error {
	return w.file.Close()
}

type Reader struct {
	repository
	file    *os.File
	decoder *json.Decoder
}

func NewFileReader(filename string) (*Reader, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Reader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (r *Reader) Read() (*model.MetricJSON, error) {
	var metric model.MetricJSON
	if err := r.decoder.Decode(&metric); err != nil {
		return nil, err
	}
	return &metric, nil
}

func (r *Reader) Close() error {
	return r.file.Close()
}
