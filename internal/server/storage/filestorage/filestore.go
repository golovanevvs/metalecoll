package filestorage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/golovanevvs/metalecoll/internal/server/model"
	"github.com/golovanevvs/metalecoll/internal/server/storage"
)

func SaveToFile(fileStoragePath string, store *storage.Storage) error {
	var str string
	var file *os.File
	_, err := os.Stat(fileStoragePath)
	if err == nil {
		fmt.Println("Open")
		file, err = os.OpenFile(fileStoragePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
	} else if os.IsNotExist(err) {
		fmt.Println("Create")
		file, err = os.Create(fileStoragePath)
		if err != nil {
			return err
		}
		fmt.Println("Create Open")
		file, err = os.OpenFile(fileStoragePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
	}

	defer file.Close()
	metrics := storage.GMs(*store)

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

func GetFromFile(fileStoragePath string, store *storage.Storage) error {
	var metric model.Metric

	file, err := os.Open(fileStoragePath)
	if err != nil {
		return err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		str := sc.Text()
		if err := json.Unmarshal([]byte(str), &metric); err != nil {
			return err
		}
		storage.SM(*store, metric)
	}
	return nil
}
