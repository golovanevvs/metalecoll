// Модуль filestorage предназначен для работы с файловым хранилищем.
package filestorage

import (
	"bufio"
	"context"
	"encoding/json"
	"os"

	"github.com/golovanevvs/metalecoll/internal/server/config"
	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/mapstorage"
	"github.com/golovanevvs/metalecoll/internal/server/model"
)

// fileStorage - тип файлового хранилища
type fileStorage struct {
	Name            string
	FileStoragePath string
}

// NewFileStorage - конструктор файлового хранилища.
func NewFileStorage(fileStoragePath string) *fileStorage {
	return &fileStorage{
		Name:            "Файловое хранилище: " + fileStoragePath,
		FileStoragePath: fileStoragePath,
	}
}

// GetNameDB возвращает название хранилища.
func (f *fileStorage) GetNameDB() string {
	return f.Name
}

// SaveMetricsToDB сохраняет метрики из map-хранилища в файл.
func (f *fileStorage) SaveMetricsToDB(ctx context.Context, c *config.Config, mapStore mapstorage.Storage) error {
	var str string
	var file *os.File
	file, err := os.OpenFile(c.Storage.FileStoragePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm|os.ModeDir)
	if err != nil {
		return err
	}

	defer file.Close()
	metrics := mapStore.GetMetricsFromMap()

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

// GetMetricsFromDB возвращает метрики из файла.
func (f *fileStorage) GetMetricsFromDB(ctx context.Context, c *config.Config) (mapstorage.Storage, error) {
	var metric model.Metric

	file, err := os.Open(c.Storage.FileStoragePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	ms := mapstorage.NewMapStorage()
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		str := sc.Text()
		if err := json.Unmarshal([]byte(str), &metric); err != nil {
			return nil, err
		}
		switch metric.MetType {
		case constants.GaugeType:
			metric.MetValue = metric.MetValue.(float64)
		case constants.CounterType:
			metric.MetValue = int64(metric.MetValue.(float64))
		}
		ms.SaveMetricToMap(metric)
	}
	return ms, nil
}

// Ping - метод-заглушка для соответствия интерфейсу.
func (f *fileStorage) Ping() error {
	return nil
}

// CloseDB - метод-заглушка для соответствия интерфейсу.
func (f *fileStorage) CloseDB() error {
	return nil
}
