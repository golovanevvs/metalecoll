package filestorage

import (
	"bufio"
	"context"
	"encoding/json"
	"os"

	"github.com/golovanevvs/metalecoll/internal/server/config"
	"github.com/golovanevvs/metalecoll/internal/server/model"
	"github.com/golovanevvs/metalecoll/internal/server/storage/mapstorage"
)

// FileStorage - тип файлового хранилища
type FileStorage struct {
	Name            string
	FileStoragePath string
}

// New - конструктор файлового хранилища
func New(c *config.Config) *FileStorage {
	return &FileStorage{
		Name:            "Файловое хранилище: " + c.Storage.FileStoragePath,
		FileStoragePath: c.Storage.FileStoragePath,
	}
}

// GetNameDB возвращает название хранилища
func (f *FileStorage) GetNameDB() string {
	return f.Name
}

// SaveMetricsToDB сохраняет метрики из map-хранилища в файл
func (f *FileStorage) SaveMetricsToDB(ctx context.Context, c *config.Config, mapStore mapstorage.Storage) error {
	var str string
	var file *os.File
	file, err := os.OpenFile(c.Storage.FileStoragePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm|os.ModeDir)
	if err != nil {
		return err
	}

	defer file.Close()
	metrics := mapStore.GetMetrics()

	for _, v := range metrics {
		enc, err := json.Marshal(v)
		if err != nil {
			return err
		}
		str += string(enc) + "\n"
	}

	_, err = file.WriteString(str)
	if err != nil {
		return err
	}

	return nil
}

// GetMetricsFromDB возвращает метрики из файла
func (f *FileStorage) GetMetricsFromDB(ctx context.Context, c *config.Config) (mapstorage.Storage, error) {
	var metric model.Metric

	file, err := os.Open(c.Storage.FileStoragePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	ms := mapstorage.New()
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		str := sc.Text()
		if err := json.Unmarshal([]byte(str), &metric); err != nil {
			return nil, err
		}
		switch metric.MetType {
		case c.MetricTypeNames.GaugeType:
			metric.MetValue = metric.MetValue.(float64)
		case c.MetricTypeNames.CounterType:
			metric.MetValue = int64(metric.MetValue.(float64))
		}
		ms.SaveMetric(metric)
	}
	return ms, nil
}

// Ping - метод-заглушка для соответствия интерфейсу
func (f *FileStorage) Ping() error {
	return nil
}
