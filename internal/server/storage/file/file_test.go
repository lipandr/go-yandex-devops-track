package file

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lipandr/go-yandex-devops-track/internal/pkg/model"
)

var (
	counterMetric = model.MetricJSON{
		ID:    "counterTest",
		MType: model.TypeCounter,
		Delta: func() *int64 { i := int64(10); return &i }(),
		Value: nil,
	}
	gaugeMetric = model.MetricJSON{
		ID:    "gaugeTest",
		MType: model.TypeGauge,
		Delta: nil,
		Value: func() *float64 { f := 5.5; return &f }(),
	}
)

func TestNewFileWriter(t *testing.T) {
	filename := "testfile.txt"
	file, err := NewFileWriter(filename)
	require.NoError(t, err, "failed to create file writer")
	defer file.Close()
	require.NotNil(t, file, "file writer is nil")
	require.NotNil(t, file.file, "file is nil")
	require.NotNil(t, file.encoder, "encoder is nil")
}

func TestNewFileReader(t *testing.T) {
	// Test correctness with valid file
	filename := "testfile.txt"
	_, err := os.Open(filename)
	require.NoError(t, err, "failed to open valid file")
	defer os.Remove(filename)
	r, err := NewFileReader(filename)
	require.NoError(t, err, "failed to create file reader")
	defer r.Close()
	require.Equal(t, r.file.Name(), filename, "expected files to be equal, but they are not")

	// Test correctness with invalid file
	_, err = NewFileReader("testfile")
	require.Error(t, err, "expected error, but got nil")
}

func TestFile(t *testing.T) {
	f, err := os.CreateTemp("", "file_write_and_read_test")
	require.NoError(t, err)
	defer func(name string) {
		_ = os.Remove(name)
	}(f.Name())

	fw, err := NewFileWriter(f.Name())
	require.NoError(t, err)
	testWrite(t, fw)
	err = fw.Close()
	require.NoError(t, err)

	fr, err := NewFileReader(f.Name())
	require.NoError(t, err)
	testRead(t, fr)
	err = fr.Close()
	require.NoError(t, err)
}

func testWrite(t *testing.T, fs *Writer) {
	require.NoError(t, fs.Write(&counterMetric))
	require.NoError(t, fs.Write(&gaugeMetric))
}

func testRead(t *testing.T, fr *Reader) {
	metric, err := fr.Read()
	require.NoError(t, err)
	require.Equal(t, &counterMetric, metric)
	metric, err = fr.Read()
	require.NoError(t, err)
	require.Equal(t, &gaugeMetric, metric)
}
